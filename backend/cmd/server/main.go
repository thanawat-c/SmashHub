package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"smashhub/internal/handlers"
	"smashhub/internal/models"
	"smashhub/internal/services"
)

func main() {
	log.Println("🏸 Starting SmashHub Backend Server...")

	// 1. Initialize PostgreSQL Connection
	db, err := models.InitDB("")
	if err != nil {
		log.Fatalf("Fatal: Cannot connect to database: %v", err)
	}

	// 2. Initialize Services
	bookingService := services.NewBookingService(db)
	slipService := services.NewSlipService(db)
	adminService := services.NewAdminService(db)
	authService := services.NewAuthService(db)
	cronWorker := services.NewCronWorker(db)

	// 3. Start Background Cron Worker (30s interval for auto-expiring bookings)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cronWorker.Start(ctx)

	// 4. Initialize HTTP Handlers
	courtHandler := handlers.NewCourtHandler(bookingService)
	bookingHandler := handlers.NewBookingHandler(bookingService, slipService, db)
	adminHandler := handlers.NewAdminHandler(adminService, bookingService, cronWorker)
	authHandler := handlers.NewAuthHandler(authService)
	sseHandler := handlers.NewSSEHandler()

	// 5. Setup Gin Router
	router := gin.Default()

	// 6. Setup CORS for LINE LIFF and Admin Frontend
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "SmashHub Backend",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// 7. Register API Routes
	api := router.Group("/api/v1")
	{
		// Real-time Event Streaming
		api.GET("/events/stream", sseHandler.StreamEvents)

		// Customer & Courts Availability
		api.GET("/courts/availability", courtHandler.GetCourtAvailability)
		api.GET("/products", courtHandler.GetProducts)

		// Authentication & User Management
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/demo-login", authHandler.DemoLogin)
			auth.GET("/me", authHandler.GetMe)
			auth.POST("/logout", authHandler.Logout)
		}

		// Booking & Slip Verification
		api.POST("/bookings", bookingHandler.CreateBooking)
		api.GET("/bookings/my-bookings", bookingHandler.GetMyBookings)
		api.GET("/bookings/:id", bookingHandler.GetBookingByID)
		api.POST("/bookings/:id/verify-slip", bookingHandler.VerifySlip)

		// Admin & POS Operations
		admin := api.Group("/admin")
		{
			admin.GET("/bookings", adminHandler.GetBookings)
			admin.PATCH("/bookings/:id/manual-approve", adminHandler.ManualApprove)
			admin.POST("/bookings/:id/add-items", adminHandler.AddItems)
			admin.PATCH("/bookings/:id/check-in", adminHandler.UpdateCheckIn)
			admin.POST("/cron/trigger", adminHandler.TriggerCron)
		}
	}

	// 8. Server Configuration & Graceful Shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("🚀 SmashHub API running on port %s (http://localhost:%s)", port, port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	cancel()
	cronWorker.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("SmashHub Server exited cleanly")
}

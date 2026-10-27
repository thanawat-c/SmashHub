package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smashhub/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	DisplayName string `json:"display_name" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Email       string `json:"email"`
	Password    string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type DemoLoginRequest struct {
	Role string `json:"role"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide all required fields: name, phone number, and password"})
		return
	}

	user, token, err := h.authService.Register(req.DisplayName, req.Phone, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Welcome to SmashHub Arena!",
		"user":    user,
		"token":   token,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide your phone number / email and password"})
		return
	}

	user, token, err := h.authService.Login(req.Identifier, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

// DemoLogin handles fast 1-click login for test users
func (h *AuthHandler) DemoLogin(c *gin.Context) {
	var req DemoLoginRequest
	_ = c.ShouldBindJSON(&req)
	if req.Role == "" {
		req.Role = "CUSTOMER"
	}

	user, token, err := h.authService.DemoLogin(req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Demo login successful",
		"user":    user,
		"token":   token,
	})
}

// GetMe returns the profile of the current user
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-Id")
	}

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	user, err := h.authService.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Logout invalidates user session
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

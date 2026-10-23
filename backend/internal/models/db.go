package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// GetDSN returns the connection string from environment or sensible defaults
func GetDSN() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "smashhub"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "smashhub_secret"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "smashhub_db"
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	timezone := os.Getenv("DB_TIMEZONE")
	if timezone == "" {
		timezone = "Asia/Bangkok"
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, user, password, dbname, sslmode, timezone)
}

// InitDB initializes PostgreSQL connection using GORM
func InitDB(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		dsn = GetDSN()
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Connection pooling
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)

	DB = db
	log.Println("Database connection established successfully")
	return db, nil
}

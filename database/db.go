package database

import (
	"fmt"
	"log"
	"os"

	"my-go-project/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init initializes the database connection and performs migrations
func Init() {
	// Database connection string from environment variables
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("POSTGRES_HOSTNAME"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_SSLMODE"),
		os.Getenv("POSTGRES_TIMEZONE"),
	)

	// Configure GORM logger for SQL query logging
	logLevel := logger.Silent
	switch os.Getenv("GORM_LOG_LEVEL") {
	case "Info":
		logLevel = logger.Info
	case "Warn":
		logLevel = logger.Warn
	case "Error":
		logLevel = logger.Error
	}
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// Connect to the database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	// Auto-migrate all registered models
	for _, model := range models.GetRegisteredModels() {
		err := DB.AutoMigrate(model)
		if err != nil {
			log.Fatalf("Failed to migrate model %T: %v", model, err)
		}
	}
}

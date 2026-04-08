package database

import (
	"log"

	"habitflow/internal/config"
	"habitflow/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance.
var DB *gorm.DB

// Init initializes the database connection and runs AutoMigrate.
func Init(cfg *config.Config) {
	var err error
	logMode := logger.Info
	if cfg.Environment == "production" {
		logMode = logger.Warn
	}

	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate all models
	err = DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Habit{},
		&models.HabitLog{},
		&models.Streak{},
		&models.PushSubscription{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	log.Println("Database connected and migrated successfully")

	// Seed the demo account
	SeedDemoAccount()
}

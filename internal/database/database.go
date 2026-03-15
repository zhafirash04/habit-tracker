package database

import (
	"log"

	"habitflow/internal/config"
	"habitflow/internal/models"

	"github.com/glebarez/sqlite"
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

	DB, err = gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Enable WAL mode for better concurrent read performance
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	sqlDB.Exec("PRAGMA journal_mode=WAL;")
	sqlDB.Exec("PRAGMA foreign_keys=ON;")

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
}

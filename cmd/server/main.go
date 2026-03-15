package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"habitflow/internal/config"
	"habitflow/internal/database"
	"habitflow/internal/router"
	"habitflow/internal/scheduler"
	"habitflow/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize database
	database.Init(cfg)

	// Setup router
	r := router.Setup(cfg)

	// Initialize and start scheduler
	db := database.DB
	pushService := services.NewPushService(db, cfg)
	scoreService := services.NewScoreService(db)
	insightService := services.NewInsightService(db)
	reportService := services.NewReportService(db, scoreService, insightService)

	sched := scheduler.New(db, pushService, reportService)
	sched.Start()

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🚀 HabitFlow server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop scheduler
	sched.Stop()

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	sqlDB, err := database.DB.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Println("✅ HabitFlow server stopped gracefully")
}

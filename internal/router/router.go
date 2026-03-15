package router

import (
	"habitflow/internal/config"
	"habitflow/internal/database"
	"habitflow/internal/handlers"
	"habitflow/internal/middleware"
	"habitflow/internal/services"

	"github.com/gin-gonic/gin"
)

// Setup configures all routes and returns the Gin engine.
func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORSMiddleware(cfg))
	r.Use(middleware.SecurityHeadersMiddleware(cfg))
	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.BodySizeLimitMiddleware(cfg.MaxBodyBytes))
	r.OPTIONS("/*path", func(c *gin.Context) { c.Status(204) })

	// Serve static PWA files
	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")
	r.Static("/icons", "./web/icons")
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/favicon.ico", "./web/icons/icon-192.png")
	r.StaticFile("/manifest.json", "./web/manifest.json")
	r.StaticFile("/sw.js", "./web/sw.js")

	// Initialize services
	db := database.DB

	authService := services.NewAuthService(db, cfg)
	habitService := services.NewHabitService(db)
	streakService := services.NewStreakService(db)
	scoreService := services.NewScoreService(db)
	insightService := services.NewInsightService(db)
	reportService := services.NewReportService(db, scoreService, insightService)
	dailyReportService := services.NewDailyReportService(db)
	pushService := services.NewPushService(db, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	habitHandler := handlers.NewHabitHandler(habitService)
	checkinHandler := handlers.NewCheckinHandler(streakService)
	reportHandler := handlers.NewReportHandler(reportService, scoreService, insightService, dailyReportService)
	pushHandler := handlers.NewPushHandler(pushService)
	healthHandler := handlers.NewHealthHandler(cfg, db)

	// API v1 routes
	api := r.Group("/api/v1")
	api.Use(middleware.RequireJSONMiddleware())
	{
		api.GET("/health", healthHandler.Liveness)
		api.GET("/health/security", healthHandler.SecuritySelfCheck)

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", middleware.LoginRateLimitMiddleware(), authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Habits
			habits := protected.Group("/habits")
			{
				habits.GET("", habitHandler.GetAll)
				habits.POST("", habitHandler.Create)
				habits.GET("/today", checkinHandler.Today)
				habits.GET("/:id", habitHandler.GetByID)
				habits.PUT("/:id", habitHandler.Update)
				habits.DELETE("/:id", habitHandler.Delete)

				// Daily check-in
				habits.POST("/:id/check", checkinHandler.Check)
				habits.DELETE("/:id/check", checkinHandler.Undo)
			}

			// Reports
			reports := protected.Group("/reports")
			{
				reports.GET("/weekly", reportHandler.Weekly)
				reports.GET("/daily", reportHandler.Daily)
				reports.GET("/score", reportHandler.Score)
				reports.GET("/insights", reportHandler.Insights)
			}

			// Push notifications
			push := protected.Group("/push")
			{
				push.GET("/vapid-key", pushHandler.VAPIDKey)
				push.POST("/subscribe", pushHandler.Subscribe)
				push.DELETE("/unsubscribe", pushHandler.Unsubscribe)
			}
		}
	}

	// SPA fallback: serve index.html for any non-API route
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/index.html")
	})

	return r
}

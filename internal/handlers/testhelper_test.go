package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"habitflow/internal/config"
	"habitflow/internal/middleware"
	"habitflow/internal/models"
	"habitflow/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const testJWTSecret = "test-secret-key-for-handlers"

func init() {
	gin.SetMode(gin.TestMode)
}

// testEnv holds all shared resources for a handler test.
type testEnv struct {
	Router *gin.Engine
	DB     *gorm.DB
	Cfg    *config.Config
}

// setupTestRouter creates a fresh in-memory DB, wires up all services/handlers/routes,
// and returns everything needed for httptest-based testing.
func setupTestRouter(t *testing.T) *testEnv {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Habit{},
		&models.HabitLog{},
		&models.Streak{},
		&models.PushSubscription{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	cfg := &config.Config{
		JWTSecret:       testJWTSecret,
		Environment:     "development",
		MaxBodyBytes:    1 << 20,
		VAPIDPublicKey:  "BEl62iUYgUivxIkv69yViEuiBIa-Ib9-SkvMeAtA3LFgDzkPs5Nqd-T3VRqGNE8raQ7n4oPk3Nl1WNbfIBmlkQ8",
		VAPIDPrivateKey: "test-vapid-private-key",
		VAPIDSubject:    "mailto:test@habitflow.app",
	}

	authService := services.NewAuthService(db, cfg)
	habitService := services.NewHabitService(db)
	streakService := services.NewStreakService(db)
	scoreService := services.NewScoreService(db)
	insightService := services.NewInsightService(db)
	reportService := services.NewReportService(db, scoreService, insightService)
	pushService := services.NewPushService(db, cfg)

	authHandler := NewAuthHandler(authService)
	habitHandler := NewHabitHandler(habitService)
	checkinHandler := NewCheckinHandler(streakService)
	dailyReportService := services.NewDailyReportService(db)
	reportHandler := NewReportHandler(reportService, scoreService, insightService, dailyReportService)
	pushHandler := NewPushHandler(pushService)
	healthHandler := NewHealthHandler(cfg, db)

	r := gin.New()

	api := r.Group("/api/v1")
	{
		api.GET("/health", healthHandler.Liveness)
		api.GET("/health/security", healthHandler.SecuritySelfCheck)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", middleware.LoginRateLimitMiddleware(), authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			habits := protected.Group("/habits")
			{
				habits.GET("", habitHandler.GetAll)
				habits.POST("", habitHandler.Create)
				habits.GET("/today", checkinHandler.Today)
				habits.GET("/:id", habitHandler.GetByID)
				habits.PUT("/:id", habitHandler.Update)
				habits.DELETE("/:id", habitHandler.Delete)
				habits.POST("/:id/check", checkinHandler.Check)
				habits.DELETE("/:id/check", checkinHandler.Undo)
			}

			reports := protected.Group("/reports")
			{
				reports.GET("/weekly", reportHandler.Weekly)
				reports.GET("/score", reportHandler.Score)
				reports.GET("/insights", reportHandler.Insights)
			}

			push := protected.Group("/push")
			{
				push.GET("/vapid-key", pushHandler.VAPIDKey)
				push.POST("/subscribe", pushHandler.Subscribe)
				push.DELETE("/unsubscribe", pushHandler.Unsubscribe)
			}
		}
	}

	return &testEnv{Router: r, DB: db, Cfg: cfg}
}

// createTestUser inserts a user with a bcrypt-hashed password and returns the model.
func createTestUser(t *testing.T, db *gorm.DB, name, email, password string) models.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	user := models.User{Name: name, Email: email, PasswordHash: string(hash)}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return user
}

// generateAccessToken creates a signed JWT access token for testing.
func generateAccessToken(t *testing.T, cfg *config.Config, userID uint, email string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"type":    "access",
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return s
}

// generateRefreshToken creates a signed JWT refresh token for testing.
func generateRefreshToken(t *testing.T, cfg *config.Config, userID uint, email string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"type":    "refresh",
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		t.Fatalf("failed to sign refresh token: %v", err)
	}
	return s
}

// generateExpiredToken creates an already-expired JWT token.
func generateExpiredToken(t *testing.T, cfg *config.Config, userID uint, email, tokenType string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"type":    tokenType,
		"exp":     time.Now().Add(-1 * time.Hour).Unix(),
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}
	return s
}

// jsonBody marshals v into an *bytes.Reader suitable for http.NewRequest.
func jsonBody(t *testing.T, v interface{}) *bytes.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal json body: %v", err)
	}
	return bytes.NewReader(b)
}

// doRequest is a convenience wrapper that creates, executes, and returns a recorded response.
func doRequest(router *gin.Engine, method, path string, body *bytes.Reader, token string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, body)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// doRequestWithCookie sends a request with optional auth header and cookie.
func doRequestWithCookie(router *gin.Engine, method, path string, body *bytes.Reader, token, cookieName, cookieValue string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, body)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if cookieName != "" {
		req.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func getCookieValue(t *testing.T, w *httptest.ResponseRecorder, name string) string {
	t.Helper()
	for _, c := range w.Result().Cookies() {
		if c.Name == name {
			return c.Value
		}
	}

	// Fallback parser for environments that expose only raw Set-Cookie headers
	for _, raw := range w.Header().Values("Set-Cookie") {
		parts := strings.Split(raw, ";")
		if len(parts) == 0 {
			continue
		}
		kv := strings.SplitN(parts[0], "=", 2)
		if len(kv) != 2 {
			continue
		}
		if strings.TrimSpace(kv[0]) == name {
			decoded, err := url.QueryUnescape(strings.TrimSpace(kv[1]))
			if err == nil {
				return decoded
			}
			return strings.TrimSpace(kv[1])
		}
	}

	t.Fatalf("cookie %s not found in response", name)
	return ""
}

// doRequestRaw sends a request with a raw io.Reader body (for malformed JSON tests).
func doRequestRaw(router *gin.Engine, method, path string, body *strings.Reader, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// parseJSON parses the recorder body into a map.
func parseJSON(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response JSON: %v\nbody: %s", err, w.Body.String())
	}
	return result
}

// assertStatus checks the HTTP status code.
func assertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if w.Code != expected {
		t.Errorf("expected status %d, got %d — body: %s", expected, w.Code, w.Body.String())
	}
}

// createTestHabit inserts a habit + streak for a user and returns the habit.
func createTestHabit(t *testing.T, db *gorm.DB, userID uint, name, category string) models.Habit {
	t.Helper()
	habit := models.Habit{UserID: userID, Name: name, Category: category, IsActive: true}
	db.Create(&habit)
	db.Create(&models.Streak{HabitID: habit.ID})
	return habit
}

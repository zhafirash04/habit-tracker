package config

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	Environment     string
	CORSOrigins     []string
	MaxBodyBytes    int64
	VAPIDPublicKey  string
	VAPIDPrivateKey string
	VAPIDSubject    string
}

// Load reads the .env file and returns a populated Config.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		JWTSecret:       getEnv("JWT_SECRET", "default-secret-change-me"),
		Environment:     strings.ToLower(getEnv("APP_ENV", "development")),
		CORSOrigins:     parseOrigins(getEnv("CORS_ALLOWED_ORIGINS", "")),
		MaxBodyBytes:    getEnvInt64("MAX_BODY_BYTES", 1<<20), // 1 MiB default
		VAPIDPublicKey:  getEnvTrimmed("VAPID_PUBLIC_KEY", ""),
		VAPIDPrivateKey: getEnvTrimmed("VAPID_PRIVATE_KEY", ""),
		VAPIDSubject:    getEnv("VAPID_SUBJECT", "mailto:admin@habitflow.app"),
	}
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL wajib diisi")
	}
	if c.MaxBodyBytes < 1024 {
		return errors.New("MAX_BODY_BYTES terlalu kecil")
	}
	if c.Environment == "production" {
		if isWeakJWTSecret(c.JWTSecret) {
			return errors.New("JWT_SECRET lemah atau default pada production")
		}
	}
	return nil
}

// IsJWTSecretStrong returns true when JWT secret passes minimum strength checks.
func (c *Config) IsJWTSecretStrong() bool {
	return !isWeakJWTSecret(c.JWTSecret)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvTrimmed(key, fallback string) string {
	value := getEnv(key, fallback)
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "\"'")
	return value
}

func getEnvInt64(key string, fallback int64) int64 {
	v := strings.TrimSpace(getEnv(key, ""))
	if v == "" {
		return fallback
	}
	var n int64
	for _, ch := range v {
		if ch < '0' || ch > '9' {
			return fallback
		}
		n = n*10 + int64(ch-'0')
	}
	if n <= 0 {
		return fallback
	}
	return n
}

func parseOrigins(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		o := strings.TrimSpace(p)
		if o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}

func isWeakJWTSecret(secret string) bool {
	s := strings.TrimSpace(secret)
	if s == "" || s == "default-secret-change-me" {
		return true
	}
	if len(s) < 32 {
		return true
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(s)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(s)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(s)
	return !(hasUpper && hasLower && hasDigit)
}

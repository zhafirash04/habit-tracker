package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"habitflow/internal/config"
	"habitflow/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const bcryptCost = 12

// AuthService handles authentication logic.
type AuthService struct {
	DB  *gorm.DB
	Cfg *config.Config
}

// NewAuthService creates a new AuthService.
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{DB: db, Cfg: cfg}
}

// RegisterInput holds the data needed to register a user.
type RegisterInput struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginInput holds the data needed to log in a user.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse is the response payload containing JWT tokens.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
}

// AuthResponse is the full response for register/login containing user data and tokens.
type AuthResponse struct {
	User   UserResponse  `json:"user"`
	Tokens TokenResponse `json:"tokens"`
}

// UserResponse is a safe user representation (no password hash).
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Register creates a new user account and returns user data with tokens.
func (s *AuthService) Register(input RegisterInput) (*AuthResponse, error) {
	// Check if email already exists
	var existing models.User
	if err := s.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcryptCost)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, errors.New("gagal membuat akun")
	}

	tokens, _, err := s.generateTokens(s.DB, user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		Tokens: *tokens,
	}, nil
}

// Login authenticates a user and returns user data with tokens.
func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	var user models.User
	if err := s.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return nil, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("email atau password salah")
	}

	tokens, _, err := s.generateTokens(s.DB, user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		Tokens: *tokens,
	}, nil
}

// RefreshToken generates a new access token from a valid refresh token.
func (s *AuthService) RefreshToken(refreshToken string) (*TokenResponse, error) {
	claims, err := s.parseRefreshClaims(refreshToken)
	if err != nil {
		return nil, err
	}

	refreshHash := hashToken(refreshToken)
	var tokens *TokenResponse
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		// Verify the refresh token exists and is still active
		var stored models.RefreshToken
		if err := tx.Where("user_id = ? AND jti = ?", claims.UserID, claims.JTI).First(&stored).Error; err != nil {
			return errors.New("refresh token tidak valid")
		}

		if stored.TokenHash != refreshHash || stored.RevokedAt != nil || time.Now().After(stored.ExpiresAt) {
			return errors.New("refresh token tidak valid")
		}

		// Verify user still exists
		var user models.User
		if err := tx.First(&user, claims.UserID).Error; err != nil {
			return errors.New("user tidak ditemukan")
		}

		newTokens, newJTI, err := s.generateTokens(tx, user.ID, user.Email)
		if err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&stored).Updates(map[string]interface{}{
			"revoked_at":      &now,
			"replaced_by_jti": newJTI,
		}).Error; err != nil {
			return errors.New("gagal merotasi refresh token")
		}

		tokens = newTokens
		return nil
	})
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// RevokeRefreshToken invalidates a refresh token so it can no longer be used.
func (s *AuthService) RevokeRefreshToken(refreshToken string) error {
	claims, err := s.parseRefreshClaims(refreshToken)
	if err != nil {
		return err
	}

	refreshHash := hashToken(refreshToken)
	var stored models.RefreshToken
	if err := s.DB.Where("user_id = ? AND jti = ?", claims.UserID, claims.JTI).First(&stored).Error; err != nil {
		return errors.New("refresh token tidak valid")
	}

	if stored.TokenHash != refreshHash {
		return errors.New("refresh token tidak valid")
	}

	if stored.RevokedAt != nil {
		return nil
	}

	now := time.Now()
	if err := s.DB.Model(&stored).Update("revoked_at", &now).Error; err != nil {
		return errors.New("gagal mencabut refresh token")
	}

	return nil
}

// generateTokens creates access and refresh JWT tokens for a user.
func (s *AuthService) generateTokens(db *gorm.DB, userID uint, email string) (*TokenResponse, string, error) {
	// Access token: 15 minutes
	accessExpiry := time.Now().Add(15 * time.Minute)
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"type":    "access",
		"exp":     accessExpiry.Unix(),
		"iat":     time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(s.Cfg.JWTSecret))
	if err != nil {
		return nil, "", errors.New("gagal membuat access token")
	}

	// Refresh token: 7 days
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
	refreshJTI := uuid.NewString()
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"type":    "refresh",
		"jti":     refreshJTI,
		"exp":     refreshExpiry.Unix(),
		"iat":     time.Now().Unix(),
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshTokenObj.SignedString([]byte(s.Cfg.JWTSecret))
	if err != nil {
		return nil, "", errors.New("gagal membuat refresh token")
	}

	refreshRecord := models.RefreshToken{
		UserID:    userID,
		JTI:       refreshJTI,
		TokenHash: hashToken(refreshString),
		ExpiresAt: refreshExpiry,
	}
	if err := db.Create(&refreshRecord).Error; err != nil {
		return nil, "", errors.New("gagal menyimpan refresh token")
	}

	return &TokenResponse{
		AccessToken:  accessString,
		RefreshToken: refreshString,
		ExpiresIn:    int64(15 * time.Minute / time.Second), // 900 seconds
	}, refreshJTI, nil
}

type refreshClaims struct {
	UserID uint
	JTI    string
}

func (s *AuthService) parseRefreshClaims(refreshToken string) (*refreshClaims, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.Cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("refresh token tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token claims tidak valid")
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "refresh" {
		return nil, errors.New("token bukan refresh token")
	}

	userIDRaw, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("user_id tidak valid dalam token")
	}

	jti, _ := claims["jti"].(string)
	if strings.TrimSpace(jti) == "" {
		return nil, errors.New("refresh token tidak valid")
	}

	return &refreshClaims{UserID: uint(userIDRaw), JTI: jti}, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

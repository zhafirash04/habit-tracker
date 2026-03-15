package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"habitflow/internal/config"
	"habitflow/internal/models"

	webpush "github.com/SherClockHolmes/webpush-go"
	"gorm.io/gorm"
)

// PushService handles Web Push notifications.
type PushService struct {
	DB  *gorm.DB
	Cfg *config.Config
}

var vapidKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
var pushKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

const maxSubscriptionsPerUser = 5

// NewPushService creates a new PushService.
func NewPushService(db *gorm.DB, cfg *config.Config) *PushService {
	return &PushService{DB: db, Cfg: cfg}
}

// SubscribeKeys holds the p256dh and auth keys from the browser PushSubscription.
type SubscribeKeys struct {
	P256dh string `json:"p256dh" binding:"required"`
	Auth   string `json:"auth" binding:"required"`
}

// SubscribeInput matches the browser PushSubscription object format.
type SubscribeInput struct {
	Endpoint string        `json:"endpoint" binding:"required"`
	Keys     SubscribeKeys `json:"keys" binding:"required"`
}

// Subscribe stores a new push subscription for a user.
func (s *PushService) Subscribe(userID uint, input SubscribeInput) (*models.PushSubscription, error) {
	if err := validateSubscribeInput(input); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidSubscription, err)
	}

	var existingCount int64
	s.DB.Model(&models.PushSubscription{}).Where("user_id = ?", userID).Count(&existingCount)

	var existingSameEndpoint int64
	s.DB.Model(&models.PushSubscription{}).
		Where("user_id = ? AND endpoint = ?", userID, input.Endpoint).
		Count(&existingSameEndpoint)

	if existingSameEndpoint == 0 && existingCount >= maxSubscriptionsPerUser {
		return nil, ErrSubscriptionLimit
	}

	// Remove existing subscription with the same endpoint to avoid duplicates
	s.DB.Where("user_id = ? AND endpoint = ?", userID, input.Endpoint).
		Delete(&models.PushSubscription{})

	sub := models.PushSubscription{
		UserID:   userID,
		Endpoint: input.Endpoint,
		P256dh:   input.Keys.P256dh,
		Auth:     input.Keys.Auth,
	}

	if err := s.DB.Create(&sub).Error; err != nil {
		return nil, err
	}

	return &sub, nil
}

// Unsubscribe removes a push subscription by endpoint.
func (s *PushService) Unsubscribe(userID uint, endpoint string) error {
	result := s.DB.Where("user_id = ? AND endpoint = ?", userID, endpoint).
		Delete(&models.PushSubscription{})
	return result.Error
}

// PushPayload is the notification payload sent to the client.
type PushPayload struct {
	Title string      `json:"title"`
	Body  string      `json:"body"`
	Icon  string      `json:"icon,omitempty"`
	Badge string      `json:"badge,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// PushData holds extra data for the notification (used by service worker).
type PushData struct {
	HabitID uint   `json:"habit_id,omitempty"`
	URL     string `json:"url,omitempty"`
}

// SendToUser sends a push notification to all subscriptions of a user.
func (s *PushService) SendToUser(userID uint, payload PushPayload) error {
	if s.Cfg.VAPIDPublicKey == "" || s.Cfg.VAPIDPrivateKey == "" {
		log.Println("VAPID keys not configured, skipping push notification")
		return nil
	}

	var subscriptions []models.PushSubscription
	s.DB.Where("user_id = ?", userID).Find(&subscriptions)

	if len(subscriptions) == 0 {
		return nil
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	for _, sub := range subscriptions {
		subscription := &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256dh,
				Auth:   sub.Auth,
			},
		}

		resp, err := webpush.SendNotification(payloadBytes, subscription, &webpush.Options{
			VAPIDPublicKey:  s.Cfg.VAPIDPublicKey,
			VAPIDPrivateKey: s.Cfg.VAPIDPrivateKey,
			Subscriber:      s.Cfg.VAPIDSubject,
		})
		if err != nil {
			log.Printf("Failed to send push to %s: %v", sub.Endpoint, err)
			if resp != nil && (resp.StatusCode == 410 || resp.StatusCode == 404) {
				// Remove invalid/expired subscriptions
				s.DB.Delete(&sub)
				log.Printf("Removed invalid subscription: %s", sub.Endpoint)
			}
			continue
		}
		resp.Body.Close()
	}

	return nil
}

// GetVAPIDPublicKey returns the VAPID public key for the frontend.
func (s *PushService) GetVAPIDPublicKey() string {
	key := sanitizeVAPIDKey(s.Cfg.VAPIDPublicKey)
	if !isLikelyValidVAPIDPublicKey(key) {
		return ""
	}
	return key
}

func validateSubscribeInput(input SubscribeInput) error {
	endpoint := strings.TrimSpace(input.Endpoint)
	if endpoint == "" || len(endpoint) > 2048 {
		return errors.New("endpoint tidak valid")
	}

	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return errors.New("endpoint harus URL https yang valid")
	}

	p256 := strings.TrimSpace(input.Keys.P256dh)
	auth := strings.TrimSpace(input.Keys.Auth)
	if len(p256) < 16 || len(p256) > 256 || !pushKeyPattern.MatchString(p256) {
		return errors.New("kunci p256dh tidak valid")
	}
	if len(auth) < 8 || len(auth) > 128 || !pushKeyPattern.MatchString(auth) {
		return errors.New("kunci auth tidak valid")
	}

	return nil
}

func sanitizeVAPIDKey(raw string) string {
	key := strings.TrimSpace(raw)
	key = strings.Trim(key, "\"'")
	return key
}

func isLikelyValidVAPIDPublicKey(key string) bool {
	if key == "" {
		return false
	}
	if strings.Contains(strings.ToLower(key), "your-vapid") {
		return false
	}
	if len(key) < 60 {
		return false
	}
	if !vapidKeyPattern.MatchString(key) {
		return false
	}
	return true
}

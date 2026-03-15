package handlers

import (
	"testing"
)

// ─── Push Subscribe Tests ─────────────────────────────────────────

func TestSubscribe_Success(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "PushUser", "push@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	body := jsonBody(t, map[string]interface{}{
		"endpoint": "https://fcm.googleapis.com/fcm/send/test-endpoint",
		"keys": map[string]string{
			"p256dh": "BF4M7u4Lh2Jd93xQ6xQ4JYq9hY3w6U8m8Mok1S6k3qAq9Tg3s3",
			"auth":   "Q2h4Qm5wY2s5U3Jm",
		},
	})

	w := doRequest(env.Router, "POST", "/api/v1/push/subscribe", body, token)
	assertStatus(t, w, 201)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

func TestSubscribe_MissingFields(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "MissP", "missp@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	// Missing keys
	body := jsonBody(t, map[string]interface{}{
		"endpoint": "https://fcm.googleapis.com/fcm/send/test",
	})

	w := doRequest(env.Router, "POST", "/api/v1/push/subscribe", body, token)
	assertStatus(t, w, 400)
}

func TestSubscribe_InvalidEndpoint(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "BadPush", "badpush@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	body := jsonBody(t, map[string]interface{}{
		"endpoint": "http://insecure-endpoint.test/push",
		"keys": map[string]string{
			"p256dh": "BF4M7u4Lh2Jd93xQ6xQ4JYq9hY3w6U8m8Mok1S6k3qAq9Tg3s3",
			"auth":   "Q2h4Qm5wY2s5U3Jm",
		},
	})

	w := doRequest(env.Router, "POST", "/api/v1/push/subscribe", body, token)
	assertStatus(t, w, 400)
}

func TestSubscribe_SubscriptionLimit(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "LimitPush", "limitpush@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	for i := 0; i < 5; i++ {
		body := jsonBody(t, map[string]interface{}{
			"endpoint": "https://fcm.googleapis.com/fcm/send/endpoint-" + string(rune('a'+i)),
			"keys": map[string]string{
				"p256dh": "BF4M7u4Lh2Jd93xQ6xQ4JYq9hY3w6U8m8Mok1S6k3qAq9Tg3s3",
				"auth":   "Q2h4Qm5wY2s5U3Jm",
			},
		})
		w := doRequest(env.Router, "POST", "/api/v1/push/subscribe", body, token)
		assertStatus(t, w, 201)
	}

	overBody := jsonBody(t, map[string]interface{}{
		"endpoint": "https://fcm.googleapis.com/fcm/send/endpoint-over",
		"keys": map[string]string{
			"p256dh": "BF4M7u4Lh2Jd93xQ6xQ4JYq9hY3w6U8m8Mok1S6k3qAq9Tg3s3",
			"auth":   "Q2h4Qm5wY2s5U3Jm",
		},
	})
	overResp := doRequest(env.Router, "POST", "/api/v1/push/subscribe", overBody, token)
	assertStatus(t, overResp, 429)
}

func TestSubscribe_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	body := jsonBody(t, map[string]interface{}{
		"endpoint": "https://fcm.googleapis.com/test",
		"keys":     map[string]string{"p256dh": "k", "auth": "a"},
	})
	w := doRequest(env.Router, "POST", "/api/v1/push/subscribe", body, "")
	assertStatus(t, w, 401)
}

// ─── Push Unsubscribe Tests ──────────────────────────────────────

func TestUnsubscribe_Success(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "UnsubU", "unsub@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	// Subscribe first
	endpoint := "https://fcm.googleapis.com/fcm/send/to-remove"
	subBody := jsonBody(t, map[string]interface{}{
		"endpoint": endpoint,
		"keys": map[string]string{
			"p256dh": "BF4M7u4Lh2Jd93xQ6xQ4JYq9hY3w6U8m8Mok1S6k3qAq9Tg3s3",
			"auth":   "Q2h4Qm5wY2s5U3Jm",
		},
	})
	doRequest(env.Router, "POST", "/api/v1/push/subscribe", subBody, token)

	// Now unsubscribe
	unsubBody := jsonBody(t, map[string]string{"endpoint": endpoint})
	w := doRequest(env.Router, "DELETE", "/api/v1/push/unsubscribe", unsubBody, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

func TestUnsubscribe_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	body := jsonBody(t, map[string]string{"endpoint": "https://example.com"})
	w := doRequest(env.Router, "DELETE", "/api/v1/push/unsubscribe", body, "")
	assertStatus(t, w, 401)
}

// ─── VAPID Key Tests ──────────────────────────────────────────────

func TestVAPIDKey_Success(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "VAPID", "vapid@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	w := doRequest(env.Router, "GET", "/api/v1/push/vapid-key", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data := resp["data"].(map[string]interface{})
	if data["vapid_public_key"] != "BEl62iUYgUivxIkv69yViEuiBIa-Ib9-SkvMeAtA3LFgDzkPs5Nqd-T3VRqGNE8raQ7n4oPk3Nl1WNbfIBmlkQ8" {
		t.Errorf("expected test vapid key, got %v", data["vapid_public_key"])
	}
}

func TestVAPIDKey_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	w := doRequest(env.Router, "GET", "/api/v1/push/vapid-key", nil, "")
	assertStatus(t, w, 401)
}

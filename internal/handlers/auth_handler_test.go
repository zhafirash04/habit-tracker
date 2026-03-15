package handlers

import (
	"strings"
	"testing"
)

// ─── Auth Handler Tests ───────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	env := setupTestRouter(t)
	body := jsonBody(t, map[string]string{
		"name": "Budi", "email": "budi@test.com", "password": "password123",
	})

	w := doRequest(env.Router, "POST", "/api/v1/auth/register", body, "")
	assertStatus(t, w, 201)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data := resp["data"].(map[string]interface{})
	if data["access_token"] == nil || data["access_token"] == "" {
		t.Error("expected access_token in response")
	}
	if _, ok := data["refresh_token"]; ok {
		t.Error("refresh_token should not be exposed in JSON response")
	}
	_ = getCookieValue(t, w, refreshCookieName)
	user := data["user"].(map[string]interface{})
	if user["name"] != "Budi" {
		t.Errorf("expected name=Budi, got %v", user["name"])
	}
	if user["email"] != "budi@test.com" {
		t.Errorf("expected email=budi@test.com, got %v", user["email"])
	}
}

func TestRegister_NameEmpty(t *testing.T) {
	env := setupTestRouter(t)
	body := jsonBody(t, map[string]string{
		"name": "", "email": "a@test.com", "password": "password123",
	})

	w := doRequest(env.Router, "POST", "/api/v1/auth/register", body, "")
	assertStatus(t, w, 400)

	resp := parseJSON(t, w)
	if resp["success"] != false {
		t.Error("expected success=false")
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	env := setupTestRouter(t)
	body := jsonBody(t, map[string]string{
		"name": "Test", "email": "not-an-email", "password": "password123",
	})

	w := doRequest(env.Router, "POST", "/api/v1/auth/register", body, "")
	assertStatus(t, w, 400)
}

func TestRegister_ShortPassword(t *testing.T) {
	env := setupTestRouter(t)
	body := jsonBody(t, map[string]string{
		"name": "Test", "email": "a@test.com", "password": "short",
	})

	w := doRequest(env.Router, "POST", "/api/v1/auth/register", body, "")
	assertStatus(t, w, 400)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	env := setupTestRouter(t)
	createTestUser(t, env.DB, "Existing", "dup@test.com", "password123")

	body := jsonBody(t, map[string]string{
		"name": "New", "email": "dup@test.com", "password": "password123",
	})

	w := doRequest(env.Router, "POST", "/api/v1/auth/register", body, "")
	assertStatus(t, w, 409)

	resp := parseJSON(t, w)
	if !strings.Contains(resp["message"].(string), "sudah") {
		t.Errorf("expected message about duplicate email, got: %v", resp["message"])
	}
}

func TestRegister_MalformedJSON(t *testing.T) {
	env := setupTestRouter(t)
	raw := strings.NewReader("{invalid json")
	w := doRequestRaw(env.Router, "POST", "/api/v1/auth/register", raw, "")
	assertStatus(t, w, 400)
}

func TestRegister_EmptyBody(t *testing.T) {
	env := setupTestRouter(t)
	w := doRequest(env.Router, "POST", "/api/v1/auth/register", nil, "")
	assertStatus(t, w, 400)
}

// ─── Login Tests ──────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	env := setupTestRouter(t)
	createTestUser(t, env.DB, "Login User", "login@test.com", "password123")

	body := jsonBody(t, map[string]string{
		"email": "login@test.com", "password": "password123",
	})

	w := doRequest(env.Router, "POST", "/api/v1/auth/login", body, "")
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data := resp["data"].(map[string]interface{})
	if data["access_token"] == nil || data["access_token"] == "" {
		t.Error("expected access_token")
	}
	if _, ok := data["refresh_token"]; ok {
		t.Error("refresh_token should not be exposed in login response")
	}
	_ = getCookieValue(t, w, refreshCookieName)
	user := data["user"].(map[string]interface{})
	if user["email"] != "login@test.com" {
		t.Errorf("expected email=login@test.com, got %v", user["email"])
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	env := setupTestRouter(t)
	createTestUser(t, env.DB, "User", "wrong@test.com", "correctpass1")

	body := jsonBody(t, map[string]string{
		"email": "wrong@test.com", "password": "wrongpass1",
	})

	w := doRequest(env.Router, "POST", "/api/v1/auth/login", body, "")
	assertStatus(t, w, 401)

	resp := parseJSON(t, w)
	if !strings.Contains(resp["message"].(string), "salah") {
		t.Errorf("expected generic error, got: %v", resp["message"])
	}
}

func TestLogin_EmailNotFound(t *testing.T) {
	env := setupTestRouter(t)

	body := jsonBody(t, map[string]string{
		"email": "noone@test.com", "password": "password123",
	})

	w := doRequest(env.Router, "POST", "/api/v1/auth/login", body, "")
	assertStatus(t, w, 401)

	resp := parseJSON(t, w)
	// Should be generic error (not revealing whether email exists)
	if !strings.Contains(resp["message"].(string), "salah") {
		t.Errorf("expected generic error, got: %v", resp["message"])
	}
}

func TestLogin_EmptyBody(t *testing.T) {
	env := setupTestRouter(t)
	w := doRequest(env.Router, "POST", "/api/v1/auth/login", nil, "")
	assertStatus(t, w, 400)
}

func TestLogin_RateLimited(t *testing.T) {
	env := setupTestRouter(t)
	createTestUser(t, env.DB, "Rate", "rate@test.com", "password123")

	body := map[string]string{"email": "rate@test.com", "password": "wrongpass1"}

	// Send 5 requests (burst size) — should all get through
	for i := 0; i < 5; i++ {
		w := doRequest(env.Router, "POST", "/api/v1/auth/login", jsonBody(t, body), "")
		if w.Code == 429 {
			t.Fatalf("request %d should not be rate limited yet", i+1)
		}
	}

	// 6th request should be rate limited
	w := doRequest(env.Router, "POST", "/api/v1/auth/login", jsonBody(t, body), "")
	assertStatus(t, w, 429)

	resp := parseJSON(t, w)
	if resp["success"] != false {
		t.Error("expected success=false on rate limit")
	}
}

// ─── Refresh Token Tests ─────────────────────────────────────────

func TestRefresh_Success(t *testing.T) {
	env := setupTestRouter(t)
	createTestUser(t, env.DB, "Refresh", "refresh@test.com", "password123")

	loginBody := jsonBody(t, map[string]string{"email": "refresh@test.com", "password": "password123"})
	loginResp := doRequest(env.Router, "POST", "/api/v1/auth/login", loginBody, "")
	assertStatus(t, loginResp, 200)
	refreshToken := getCookieValue(t, loginResp, refreshCookieName)

	w := doRequestWithCookie(env.Router, "POST", "/api/v1/auth/refresh", nil, "", refreshCookieName, refreshToken)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data := resp["data"].(map[string]interface{})
	if data["access_token"] == nil || data["access_token"] == "" {
		t.Error("expected new access_token")
	}
	if _, ok := data["refresh_token"]; ok {
		t.Error("refresh_token should not be exposed in refresh response")
	}
	_ = getCookieValue(t, w, refreshCookieName)
}

func TestRefresh_ExpiredToken(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Expired", "expired@test.com", "password123")
	expired := generateExpiredToken(t, env.Cfg, user.ID, user.Email, "refresh")

	w := doRequest(env.Router, "POST", "/api/v1/auth/refresh", nil, expired)
	assertStatus(t, w, 401)
}

func TestRefresh_InvalidToken(t *testing.T) {
	env := setupTestRouter(t)

	w := doRequest(env.Router, "POST", "/api/v1/auth/refresh", nil, "completely-invalid-token")
	assertStatus(t, w, 401)
}

func TestRefresh_RotationRejectsOldToken(t *testing.T) {
	env := setupTestRouter(t)
	createTestUser(t, env.DB, "Rotate", "rotate@test.com", "password123")

	loginBody := jsonBody(t, map[string]string{"email": "rotate@test.com", "password": "password123"})
	loginResp := doRequest(env.Router, "POST", "/api/v1/auth/login", loginBody, "")
	assertStatus(t, loginResp, 200)
	oldRefresh := getCookieValue(t, loginResp, refreshCookieName)

	firstRefresh := doRequestWithCookie(env.Router, "POST", "/api/v1/auth/refresh", nil, "", refreshCookieName, oldRefresh)
	assertStatus(t, firstRefresh, 200)

	secondRefresh := doRequestWithCookie(env.Router, "POST", "/api/v1/auth/refresh", nil, "", refreshCookieName, oldRefresh)
	assertStatus(t, secondRefresh, 401)
}

func TestLogout_RevokesRefreshToken(t *testing.T) {
	env := setupTestRouter(t)
	createTestUser(t, env.DB, "Logout", "logout@test.com", "password123")

	loginBody := jsonBody(t, map[string]string{"email": "logout@test.com", "password": "password123"})
	loginResp := doRequest(env.Router, "POST", "/api/v1/auth/login", loginBody, "")
	assertStatus(t, loginResp, 200)
	refresh := getCookieValue(t, loginResp, refreshCookieName)

	logoutResp := doRequestWithCookie(env.Router, "POST", "/api/v1/auth/logout", nil, "", refreshCookieName, refresh)
	assertStatus(t, logoutResp, 200)
	clearedCookie := getCookieValue(t, logoutResp, refreshCookieName)
	if clearedCookie != "" {
		t.Error("refresh cookie should be cleared on logout")
	}

	refreshResp := doRequestWithCookie(env.Router, "POST", "/api/v1/auth/refresh", nil, "", refreshCookieName, refresh)
	assertStatus(t, refreshResp, 401)
}

func TestRefresh_NoHeader(t *testing.T) {
	env := setupTestRouter(t)

	w := doRequest(env.Router, "POST", "/api/v1/auth/refresh", nil, "")
	assertStatus(t, w, 400)
}

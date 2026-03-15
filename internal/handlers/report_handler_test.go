package handlers

import (
	"testing"
)

// ─── Weekly Report Tests ──────────────────────────────────────────

func TestWeekly_Default(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Weekly", "weekly@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	w := doRequest(env.Router, "GET", "/api/v1/reports/weekly", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	if resp["data"] == nil {
		t.Error("expected data to be present")
	}
}

func TestWeekly_WithDateRange(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Range", "range@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	w := doRequest(env.Router, "GET", "/api/v1/reports/weekly?start=2025-01-01&end=2025-01-07", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

func TestWeekly_NoActivity(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "NoAct", "noact@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	// User has no habits, no logs — report should still work with zero values
	w := doRequest(env.Router, "GET", "/api/v1/reports/weekly", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	data := resp["data"].(map[string]interface{})
	if rate, ok := data["completion_rate"].(float64); ok && rate != 0 {
		t.Errorf("expected 0%% completion for no activity, got %v", rate)
	}
}

func TestWeekly_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	w := doRequest(env.Router, "GET", "/api/v1/reports/weekly", nil, "")
	assertStatus(t, w, 401)
}

// ─── Score Tests ──────────────────────────────────────────────────

func TestScore_Default(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "ScoreU", "score@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	w := doRequest(env.Router, "GET", "/api/v1/reports/score", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

func TestScore_CustomDays(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Score30", "score30@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	w := doRequest(env.Router, "GET", "/api/v1/reports/score?days=30", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

func TestScore_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	w := doRequest(env.Router, "GET", "/api/v1/reports/score", nil, "")
	assertStatus(t, w, 401)
}

// ─── Insights Tests ───────────────────────────────────────────────

func TestInsights_Success(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Insight", "insight@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	w := doRequest(env.Router, "GET", "/api/v1/reports/insights", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

func TestInsights_NewUser(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "NewI", "newi@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	// Brand-new user with no habits should get valid insights (possibly empty)
	w := doRequest(env.Router, "GET", "/api/v1/reports/insights", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["data"] == nil {
		t.Error("expected data to be present even for new user")
	}
}

func TestInsights_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	w := doRequest(env.Router, "GET", "/api/v1/reports/insights", nil, "")
	assertStatus(t, w, 401)
}

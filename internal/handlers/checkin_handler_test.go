package handlers

import (
	"fmt"
	"testing"
)

// ─── Checkin Handler Tests ────────────────────────────────────────

func TestCheck_Success(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Checker", "check@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)
	habit := createTestHabit(t, env.DB, user.ID, "Morning Run", "fitness")

	w := doRequest(env.Router, "POST", fmt.Sprintf("/api/v1/habits/%d/check", habit.ID), nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data := resp["data"].(map[string]interface{})
	if data["is_done"] != true {
		t.Error("expected is_done=true")
	}
	if data["current_streak"] != float64(1) {
		t.Errorf("expected current_streak=1, got %v", data["current_streak"])
	}
}

func TestCheck_WithNote(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Noter", "note@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)
	habit := createTestHabit(t, env.DB, user.ID, "Read", "learning")

	body := jsonBody(t, map[string]string{"note": "Finished chapter 5"})
	w := doRequest(env.Router, "POST", fmt.Sprintf("/api/v1/habits/%d/check", habit.ID), body, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	data := resp["data"].(map[string]interface{})
	if data["note"] != "Finished chapter 5" {
		t.Errorf("expected note, got %v", data["note"])
	}
}

func TestCheck_Duplicate(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "DupCheck", "dupch@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)
	habit := createTestHabit(t, env.DB, user.ID, "Meditate", "wellness")

	// First check-in — OK
	doRequest(env.Router, "POST", fmt.Sprintf("/api/v1/habits/%d/check", habit.ID), nil, token)

	// Second check-in — conflict
	w := doRequest(env.Router, "POST", fmt.Sprintf("/api/v1/habits/%d/check", habit.ID), nil, token)
	assertStatus(t, w, 409)

	resp := parseJSON(t, w)
	if resp["success"] != false {
		t.Error("expected success=false on duplicate")
	}
}

func TestCheck_HabitNotFound(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "NF", "nfch@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	w := doRequest(env.Router, "POST", "/api/v1/habits/9999/check", nil, token)
	assertStatus(t, w, 404)
}

func TestCheck_OtherUsersHabit(t *testing.T) {
	env := setupTestRouter(t)
	owner := createTestUser(t, env.DB, "Owner", "ownch@test.com", "password123")
	other := createTestUser(t, env.DB, "Other", "othch@test.com", "password123")
	tokenOther := generateAccessToken(t, env.Cfg, other.ID, other.Email)

	habit := createTestHabit(t, env.DB, owner.ID, "Private", "health")

	w := doRequest(env.Router, "POST", fmt.Sprintf("/api/v1/habits/%d/check", habit.ID), nil, tokenOther)
	assertStatus(t, w, 404) // ownership check → 404
}

func TestCheck_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	w := doRequest(env.Router, "POST", "/api/v1/habits/1/check", nil, "")
	assertStatus(t, w, 401)
}

// ─── Undo Tests ───────────────────────────────────────────────────

func TestUndo_Success(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Undoer", "undo@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)
	habit := createTestHabit(t, env.DB, user.ID, "Gym", "fitness")

	// Check in first
	doRequest(env.Router, "POST", fmt.Sprintf("/api/v1/habits/%d/check", habit.ID), nil, token)

	// Undo
	w := doRequest(env.Router, "DELETE", fmt.Sprintf("/api/v1/habits/%d/check", habit.ID), nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

func TestUndo_NoCheckinToday(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "NoUndo", "noundo@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)
	habit := createTestHabit(t, env.DB, user.ID, "Walk", "fitness")

	w := doRequest(env.Router, "DELETE", fmt.Sprintf("/api/v1/habits/%d/check", habit.ID), nil, token)
	assertStatus(t, w, 404)
}

func TestUndo_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	w := doRequest(env.Router, "DELETE", "/api/v1/habits/1/check", nil, "")
	assertStatus(t, w, 401)
}

// ─── Today Status Tests ──────────────────────────────────────────

func TestToday_Success(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "TodayU", "today@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	h1 := createTestHabit(t, env.DB, user.ID, "Run", "fitness")
	createTestHabit(t, env.DB, user.ID, "Read", "learning")

	// Check in only the first habit
	doRequest(env.Router, "POST", fmt.Sprintf("/api/v1/habits/%d/check", h1.ID), nil, token)

	w := doRequest(env.Router, "GET", "/api/v1/habits/today", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	data := resp["data"].(map[string]interface{})
	habits := data["habits"].([]interface{})
	if len(habits) != 2 {
		t.Errorf("expected 2 habits in today status, got %d", len(habits))
	}

	// Verify one is done, one is not
	doneCount := 0
	for _, h := range habits {
		hm := h.(map[string]interface{})
		if hm["is_done_today"] == true {
			doneCount++
		}
	}
	if doneCount != 1 {
		t.Errorf("expected 1 done habit, got %d", doneCount)
	}
}

func TestToday_EmptyHabits(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Empty", "emptytoday@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	w := doRequest(env.Router, "GET", "/api/v1/habits/today", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	data := resp["data"].(map[string]interface{})
	habits := data["habits"].([]interface{})
	if len(habits) != 0 {
		t.Errorf("expected 0 habits, got %d", len(habits))
	}
}

func TestToday_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	w := doRequest(env.Router, "GET", "/api/v1/habits/today", nil, "")
	assertStatus(t, w, 401)
}

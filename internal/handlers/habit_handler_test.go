package handlers

import (
	"fmt"
	"strings"
	"testing"
)

// ─── Habit Handler Tests ──────────────────────────────────────────

func TestGetAllHabits_Success(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Habit User", "habit@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	createTestHabit(t, env.DB, user.ID, "Minum Air", "health")
	createTestHabit(t, env.DB, user.ID, "Olahraga", "fitness")

	w := doRequest(env.Router, "GET", "/api/v1/habits", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data := resp["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("expected 2 habits, got %d", len(data))
	}
}

func TestGetAllHabits_Empty(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "New User", "new@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	w := doRequest(env.Router, "GET", "/api/v1/habits", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	data := resp["data"].([]interface{})
	if len(data) != 0 {
		t.Errorf("expected 0 habits, got %d", len(data))
	}
}

func TestGetAllHabits_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	w := doRequest(env.Router, "GET", "/api/v1/habits", nil, "")
	assertStatus(t, w, 401)
}

func TestGetAllHabits_UserIsolation(t *testing.T) {
	env := setupTestRouter(t)
	userA := createTestUser(t, env.DB, "UserA", "a@test.com", "password123")
	userB := createTestUser(t, env.DB, "UserB", "b@test.com", "password123")
	tokenB := generateAccessToken(t, env.Cfg, userB.ID, userB.Email)

	createTestHabit(t, env.DB, userA.ID, "A's Habit", "health")

	w := doRequest(env.Router, "GET", "/api/v1/habits", nil, tokenB)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	data := resp["data"].([]interface{})
	if len(data) != 0 {
		t.Error("user B should not see user A's habits")
	}
}

func TestCreateHabit_Success(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Creator", "create@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	body := jsonBody(t, map[string]interface{}{
		"name": "Baca Buku", "category": "learning", "notify_time": "07:00",
	})

	w := doRequest(env.Router, "POST", "/api/v1/habits", body, token)
	assertStatus(t, w, 201)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data := resp["data"].(map[string]interface{})
	if data["name"] != "Baca Buku" {
		t.Errorf("expected name=Baca Buku, got %v", data["name"])
	}
	if data["category"] != "learning" {
		t.Errorf("expected category=learning, got %v", data["category"])
	}
	if data["is_active"] != true {
		t.Error("expected is_active=true")
	}
	if data["current_streak"] != float64(0) {
		t.Error("expected initial current_streak=0")
	}
}

func TestCreateHabit_WithoutOptionalFields(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Minimal", "min@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	body := jsonBody(t, map[string]string{"name": "Walk"})
	w := doRequest(env.Router, "POST", "/api/v1/habits", body, token)
	assertStatus(t, w, 201)

	resp := parseJSON(t, w)
	data := resp["data"].(map[string]interface{})
	// Default category should be "general"
	if data["category"] != "general" {
		t.Errorf("expected default category=general, got %v", data["category"])
	}
}

func TestCreateHabit_NameEmpty(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Empty", "empty@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	body := jsonBody(t, map[string]string{"name": ""})
	w := doRequest(env.Router, "POST", "/api/v1/habits", body, token)
	assertStatus(t, w, 400)
}

func TestCreateHabit_NameTooLong(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Long", "long@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	longName := strings.Repeat("a", 101)
	body := jsonBody(t, map[string]string{"name": longName})
	w := doRequest(env.Router, "POST", "/api/v1/habits", body, token)
	assertStatus(t, w, 400)
}

func TestCreateHabit_InvalidNotifyTime(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Time", "time@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	body := jsonBody(t, map[string]string{"name": "Test", "notify_time": "25:99"})
	w := doRequest(env.Router, "POST", "/api/v1/habits", body, token)
	assertStatus(t, w, 400)
}

func TestCreateHabit_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	body := jsonBody(t, map[string]string{"name": "Test"})
	w := doRequest(env.Router, "POST", "/api/v1/habits", body, "")
	assertStatus(t, w, 401)
}

func TestUpdateHabit_Success(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Updater", "upd@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)
	habit := createTestHabit(t, env.DB, user.ID, "Old Name", "health")

	body := jsonBody(t, map[string]string{"name": "New Name"})
	w := doRequest(env.Router, "PUT", fmt.Sprintf("/api/v1/habits/%d", habit.ID), body, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	data := resp["data"].(map[string]interface{})
	if data["name"] != "New Name" {
		t.Errorf("expected name=New Name, got %v", data["name"])
	}
}

func TestUpdateHabit_NotifyTime(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "TimerU", "tu@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)
	habit := createTestHabit(t, env.DB, user.ID, "Habit", "health")

	body := jsonBody(t, map[string]string{"notify_time": "08:30"})
	w := doRequest(env.Router, "PUT", fmt.Sprintf("/api/v1/habits/%d", habit.ID), body, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	data := resp["data"].(map[string]interface{})
	if data["notify_time"] != "08:30" {
		t.Errorf("expected notify_time=08:30, got %v", data["notify_time"])
	}
}

func TestUpdateHabit_NotFound(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "NF", "nf@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	body := jsonBody(t, map[string]string{"name": "X"})
	w := doRequest(env.Router, "PUT", "/api/v1/habits/9999", body, token)
	assertStatus(t, w, 404)
}

func TestUpdateHabit_OtherUsersHabit(t *testing.T) {
	env := setupTestRouter(t)
	owner := createTestUser(t, env.DB, "Owner", "owner@test.com", "password123")
	other := createTestUser(t, env.DB, "Other", "other@test.com", "password123")
	tokenOther := generateAccessToken(t, env.Cfg, other.ID, other.Email)

	habit := createTestHabit(t, env.DB, owner.ID, "Owners Habit", "health")

	body := jsonBody(t, map[string]string{"name": "Hacked"})
	w := doRequest(env.Router, "PUT", fmt.Sprintf("/api/v1/habits/%d", habit.ID), body, tokenOther)
	// Owner mismatch → returns 404 (doesn't reveal existence)
	assertStatus(t, w, 404)
}

func TestUpdateHabit_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	body := jsonBody(t, map[string]string{"name": "X"})
	w := doRequest(env.Router, "PUT", "/api/v1/habits/1", body, "")
	assertStatus(t, w, 401)
}

func TestDeleteHabit_Success(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "Deleter", "del@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)
	habit := createTestHabit(t, env.DB, user.ID, "To Delete", "health")

	w := doRequest(env.Router, "DELETE", fmt.Sprintf("/api/v1/habits/%d", habit.ID), nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

func TestDeleteHabit_VerifyNotInGetAll(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "VerDel", "vd@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)
	habit := createTestHabit(t, env.DB, user.ID, "Disappear", "health")
	createTestHabit(t, env.DB, user.ID, "Stay", "health")

	// Delete one habit
	doRequest(env.Router, "DELETE", fmt.Sprintf("/api/v1/habits/%d", habit.ID), nil, token)

	// GET should only return the remaining habit
	w := doRequest(env.Router, "GET", "/api/v1/habits", nil, token)
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	data := resp["data"].([]interface{})
	if len(data) != 1 {
		t.Errorf("expected 1 remaining habit, got %d", len(data))
	}
	first := data[0].(map[string]interface{})
	if first["name"] != "Stay" {
		t.Errorf("expected remaining habit=Stay, got %v", first["name"])
	}
}

func TestDeleteHabit_NotFound(t *testing.T) {
	env := setupTestRouter(t)
	user := createTestUser(t, env.DB, "DNF", "dnf@test.com", "password123")
	token := generateAccessToken(t, env.Cfg, user.ID, user.Email)

	w := doRequest(env.Router, "DELETE", "/api/v1/habits/9999", nil, token)
	assertStatus(t, w, 404)
}

func TestDeleteHabit_OtherUsersHabit(t *testing.T) {
	env := setupTestRouter(t)
	owner := createTestUser(t, env.DB, "Own", "own@test.com", "password123")
	other := createTestUser(t, env.DB, "Oth", "oth@test.com", "password123")
	tokenOther := generateAccessToken(t, env.Cfg, other.ID, other.Email)

	habit := createTestHabit(t, env.DB, owner.ID, "Protected", "health")

	w := doRequest(env.Router, "DELETE", fmt.Sprintf("/api/v1/habits/%d", habit.ID), nil, tokenOther)
	assertStatus(t, w, 404)
}

func TestDeleteHabit_Unauthorized(t *testing.T) {
	env := setupTestRouter(t)
	w := doRequest(env.Router, "DELETE", "/api/v1/habits/1", nil, "")
	assertStatus(t, w, 401)
}

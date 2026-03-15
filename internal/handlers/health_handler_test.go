package handlers

import "testing"

func TestHealth_Liveness(t *testing.T) {
	env := setupTestRouter(t)

	w := doRequest(env.Router, "GET", "/api/v1/health", nil, "")
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data := resp["data"].(map[string]interface{})
	if data["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", data["status"])
	}
}

func TestHealth_SecuritySelfCheck(t *testing.T) {
	env := setupTestRouter(t)

	w := doRequest(env.Router, "GET", "/api/v1/health/security", nil, "")
	assertStatus(t, w, 200)

	resp := parseJSON(t, w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data := resp["data"].(map[string]interface{})
	checks := data["checks"].(map[string]interface{})
	if checks["database_connected"] != true {
		t.Error("expected database_connected=true")
	}
	if checks["refresh_token_http_only"] != true {
		t.Error("expected refresh_token_http_only=true")
	}
	if checks["api_response_cache_disabled"] != true {
		t.Error("expected api_response_cache_disabled=true")
	}
}

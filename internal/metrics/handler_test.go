package metrics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/driftwatch/internal/metrics"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	c := metrics.New()
	c.RecordCheck()
	c.RecordDrift([]string{"web"})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metrics.Handler(c)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("unexpected Content-Type: %s", ct)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if body["checks_total"].(float64) != 1 {
		t.Fatalf("expected checks_total=1, got %v", body["checks_total"])
	}
	if body["drift_total"].(float64) != 1 {
		t.Fatalf("expected drift_total=1, got %v", body["drift_total"])
	}
}

func TestHandler_EmptyContainersDriftedIsArray(t *testing.T) {
	c := metrics.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metrics.Handler(c)(rec, req)

	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if _, ok := body["containers_drifted"].([]interface{}); !ok {
		t.Fatalf("containers_drifted should be an array, got %T", body["containers_drifted"])
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	c := metrics.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/metrics", nil)
	metrics.Handler(c)(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

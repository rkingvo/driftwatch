package alerter_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/alerter"
)

func TestWebhookNotifier_PostsJSON(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wn := alerter.NewWebhookNotifier(ts.URL, nil)
	alrt := alerter.Alert{
		Timestamp: time.Now().UTC(),
		Level:     alerter.LevelWarn,
		Report:    hasDrift(),
	}
	if err := wn.Notify(alrt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["level"] != string(alerter.LevelWarn) {
		t.Errorf("expected level WARN, got %v", received["level"])
	}
	containers, ok := received["drifted_containers"].([]interface{})
	if !ok || len(containers) == 0 {
		t.Errorf("expected drifted_containers in payload, got %v", received)
	}
}

func TestWebhookNotifier_Non2xx_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	wn := alerter.NewWebhookNotifier(ts.URL, nil)
	alrt := alerter.Alert{
		Timestamp: time.Now().UTC(),
		Level:     alerter.LevelWarn,
		Report:    hasDrift(),
	}
	if err := wn.Notify(alrt); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

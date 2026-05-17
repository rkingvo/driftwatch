package alerter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier posts a JSON-encoded Alert to an HTTP endpoint.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

// NewWebhookNotifier creates a WebhookNotifier for the given URL.
// A default 10-second timeout is applied when client is nil.
func NewWebhookNotifier(url string, client *http.Client) *WebhookNotifier {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &WebhookNotifier{url: url, client: client}
}

// webhookPayload is the JSON body sent to the webhook endpoint.
type webhookPayload struct {
	Timestamp string            `json:"timestamp"`
	Level     Level             `json:"level"`
	Drifted   []string          `json:"drifted_containers"`
}

// Notify serialises the alert and POSTs it to the configured URL.
func (w *WebhookNotifier) Notify(a Alert) error {
	var names []string
	for _, c := range a.Report.Containers {
		names = append(names, c.Name)
	}
	payload := webhookPayload{
		Timestamp: a.Timestamp.Format(time.RFC3339),
		Level:     a.Level,
		Drifted:   names,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}

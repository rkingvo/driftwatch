package metrics

import (
	"encoding/json"
	"net/http"
	"time"
)

// jsonSnapshot is the wire format for the /metrics endpoint.
type jsonSnapshot struct {
	ChecksTotal       int64     `json:"checks_total"`
	DriftTotal        int64     `json:"drift_total"`
	LastCheckAt       time.Time `json:"last_check_at,omitempty"`
	LastDriftAt       time.Time `json:"last_drift_at,omitempty"`
	ContainersDrifted []string  `json:"containers_drifted"`
}

// Handler returns an http.HandlerFunc that serialises the current Snapshot
// as JSON. It is safe to register on any ServeMux.
func Handler(c *Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s := c.Snapshot()
		out := jsonSnapshot{
			ChecksTotal:       s.ChecksTotal,
			DriftTotal:        s.DriftTotal,
			LastCheckAt:       s.LastCheckAt,
			LastDriftAt:       s.LastDriftAt,
			ContainersDrifted: s.ContainersDrifted,
		}
		if out.ContainersDrifted == nil {
			out.ContainersDrifted = []string{}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(out); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	}
}

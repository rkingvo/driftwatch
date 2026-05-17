// Package healthz provides a simple HTTP health-check endpoint for the
// driftwatch daemon. It exposes /healthz (liveness) and /readyz (readiness)
// handlers that can be mounted on any http.ServeMux.
package healthz

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// Status holds the payload returned by both endpoints.
type Status struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// Handler exposes /healthz and /readyz routes.
type Handler struct {
	// ready is flipped to 1 once the daemon has completed its first check
	// cycle and is considered ready to serve traffic.
	ready atomic.Int32
	started time.Time
}

// New creates a Handler. Call MarkReady once the daemon is initialised.
func New() *Handler {
	return &Handler{started: time.Now()}
}

// MarkReady signals that the daemon has finished its first check cycle.
func (h *Handler) MarkReady() {
	h.ready.Store(1)
}

// Register mounts /healthz and /readyz on mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.liveness)
	mux.HandleFunc("/readyz", h.readiness)
}

func (h *Handler) liveness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, Status{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) readiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.ready.Load() == 0 {
		writeJSON(w, http.StatusServiceUnavailable, Status{
			Status:    "not ready",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}
	writeJSON(w, http.StatusOK, Status{
		Status:    "ready",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

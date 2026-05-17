// Package server provides the HTTP server that exposes driftwatch's
// operational endpoints.
//
// Endpoints:
//
//	/metrics        – JSON snapshot of drift metrics (GET)
//	/healthz/live   – liveness probe (GET)
//	/healthz/ready  – readiness probe (GET)
//
// Usage:
//
//	cfg := server.DefaultConfig()
//	srv := server.New(cfg, metricsCollector, healthHandler)
//	if err := srv.Start(ctx); err != nil {
//		log.Fatal(err)
//	}
package server

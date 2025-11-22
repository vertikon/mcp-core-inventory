package http

import (
	"net/http"

	"github.com/vertikon/mcp-core-inventory/internal/app"
	"github.com/vertikon/mcp-core-inventory/internal/observability"
	"go.uber.org/zap"
)

// Router sets up HTTP routes for the inventory API
type Router struct {
	handlers *ReservationHandlers
	logger   *zap.Logger
}

// NewRouter creates a new HTTP router
func NewRouter(
	reserveUseCase *app.ReserveStockUseCase,
	confirmUseCase *app.ConfirmReservationUseCase,
	releaseUseCase *app.ReleaseReservationUseCase,
	adjustUseCase *app.AdjustStockUseCase,
	queryUseCase *app.QueryAvailableUseCase,
	logger *zap.Logger,
	metrics interface{}, // Can be *observability.InventoryMetrics or nil
) *Router {
	var inventoryMetrics *observability.InventoryMetrics
	if m, ok := metrics.(*observability.InventoryMetrics); ok {
		inventoryMetrics = m
	}

	handlers := NewReservationHandlers(
		reserveUseCase,
		confirmUseCase,
		releaseUseCase,
		adjustUseCase,
		queryUseCase,
		logger,
		inventoryMetrics,
	)

	return &Router{
		handlers: handlers,
		logger:   logger,
	}
}

// SetupRoutes configures all HTTP routes
func (r *Router) SetupRoutes(mux *http.ServeMux) {
	// v1 endpoints
	mux.HandleFunc("/v1/reserve", r.handlers.ReserveStock)
	mux.HandleFunc("/v1/confirm", r.handlers.ConfirmReservation)
	mux.HandleFunc("/v1/release", r.handlers.ReleaseReservation)
	mux.HandleFunc("/v1/adjust", r.handlers.AdjustStock)
	mux.HandleFunc("/v1/available", r.handlers.QueryAvailable)

	// v2 endpoints (same functionality, with deprecation headers)
	mux.HandleFunc("/v2/reserve", r.withDeprecationHeaders(r.handlers.ReserveStock))
	mux.HandleFunc("/v2/confirm", r.withDeprecationHeaders(r.handlers.ConfirmReservation))
	mux.HandleFunc("/v2/release", r.withDeprecationHeaders(r.handlers.ReleaseReservation))
	mux.HandleFunc("/v2/adjust", r.withDeprecationHeaders(r.handlers.AdjustStock))
	mux.HandleFunc("/v2/available", r.withDeprecationHeaders(r.handlers.QueryAvailable))

	// Health check
	mux.HandleFunc("/health", r.healthCheck)
}

// withDeprecationHeaders adds deprecation headers for v2 endpoints
func (r *Router) withDeprecationHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Deprecation", "true")
		w.Header().Set("Sunset", "2026-12-31")
		w.Header().Set("Link", "</v1/docs>; rel=\"deprecation\"")
		handler(w, req)
	}
}

// healthCheck handles health check requests
func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/vertikon/mcp-core-inventory/internal/app"
	"github.com/vertikon/mcp-core-inventory/internal/observability"
	"go.uber.org/zap"
)

// ReservationHandlers handles HTTP requests for stock reservations
type ReservationHandlers struct {
	reserveUseCase *app.ReserveStockUseCase
	confirmUseCase *app.ConfirmReservationUseCase
	releaseUseCase *app.ReleaseReservationUseCase
	adjustUseCase  *app.AdjustStockUseCase
	queryUseCase   *app.QueryAvailableUseCase
	logger         *zap.Logger
	metrics        *observability.InventoryMetrics
}

// NewReservationHandlers creates new reservation handlers
func NewReservationHandlers(
	reserveUseCase *app.ReserveStockUseCase,
	confirmUseCase *app.ConfirmReservationUseCase,
	releaseUseCase *app.ReleaseReservationUseCase,
	adjustUseCase *app.AdjustStockUseCase,
	queryUseCase *app.QueryAvailableUseCase,
	logger *zap.Logger,
	metrics *observability.InventoryMetrics,
) *ReservationHandlers {
	return &ReservationHandlers{
		reserveUseCase: reserveUseCase,
		confirmUseCase: confirmUseCase,
		releaseUseCase: releaseUseCase,
		adjustUseCase:  adjustUseCase,
		queryUseCase:   queryUseCase,
		logger:         logger,
		metrics:        metrics,
	}
}

// ReserveStock handles POST /reserve
func (h *ReservationHandlers) ReserveStock(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SKU            string        `json:"sku"`
		Location       string        `json:"location"`
		Quantity       int64         `json:"quantity"`
		IdempotencyKey string        `json:"idempotency_key,omitempty"`
		TTL            time.Duration `json:"ttl,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		if h.metrics != nil {
			h.metrics.RecordReservationError(req.SKU, req.Location, "invalid_request")
		}
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	appReq := app.ReserveStockRequest{
		SKU:            req.SKU,
		Location:       req.Location,
		Quantity:       req.Quantity,
		IdempotencyKey: req.IdempotencyKey,
		TTL:            req.TTL,
	}

	resp, err := h.reserveUseCase.Execute(r.Context(), appReq)
	duration := time.Since(startTime).Seconds()

	if err != nil {
		h.logger.Error("failed to reserve stock", zap.Error(err))

		// Record error metrics
		if h.metrics != nil {
			errorType := "unknown"
			if strings.Contains(err.Error(), "insufficient") {
				errorType = "insufficient_stock"
			} else if strings.Contains(err.Error(), "lock") {
				errorType = "lock_failed"
			} else if strings.Contains(err.Error(), "race") {
				errorType = "race_condition"
				h.metrics.RecordRaceCondition()
			}
			h.metrics.RecordReservationError(req.SKU, req.Location, errorType)
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Record success metrics
	if h.metrics != nil {
		h.metrics.RecordReservation(req.SKU, req.Location, "success", duration)
		if req.IdempotencyKey != "" {
			// Check if this was an idempotency hit (existing reservation)
			// This would need to be tracked in the use case, but for now we assume new
			h.metrics.RecordIdempotencyMiss()
		} else {
			h.metrics.RecordIdempotencyMiss()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ConfirmReservation handles POST /confirm
func (h *ReservationHandlers) ConfirmReservation(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ReservationID string `json:"reservation_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	appReq := app.ConfirmReservationRequest{
		ReservationID: req.ReservationID,
	}

	resp, err := h.confirmUseCase.Execute(r.Context(), appReq)
	duration := time.Since(startTime).Seconds()

	if err != nil {
		h.logger.Error("failed to confirm reservation", zap.Error(err))

		// Record error metrics
		if h.metrics != nil {
			errorType := "unknown"
			if strings.Contains(err.Error(), "not found") {
				errorType = "reservation_not_found"
			} else if strings.Contains(err.Error(), "invalid") {
				errorType = "invalid_reservation"
			}
			// Use empty strings when resp is nil
			sku := ""
			location := ""
			if resp != nil {
				sku = resp.SKU
				location = resp.Location
			}
			h.metrics.RecordConfirmationError(sku, location, errorType)
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Record success metrics
	if h.metrics != nil {
		h.metrics.RecordConfirmation(resp.SKU, resp.Location, "success", duration)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ReleaseReservation handles POST /release
func (h *ReservationHandlers) ReleaseReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ReservationID string `json:"reservation_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	appReq := app.ReleaseReservationRequest{
		ReservationID: req.ReservationID,
	}

	resp, err := h.releaseUseCase.Execute(r.Context(), appReq)

	if err != nil {
		h.logger.Error("failed to release reservation", zap.Error(err))

		// Record error metrics
		if h.metrics != nil {
			errorType := "unknown"
			if strings.Contains(err.Error(), "not found") {
				errorType = "reservation_not_found"
			} else if strings.Contains(err.Error(), "invalid") {
				errorType = "invalid_reservation"
			}
			// Use empty strings when resp is nil
			sku := ""
			location := ""
			if resp != nil {
				sku = resp.SKU
				location = resp.Location
			}
			h.metrics.RecordReleaseError(sku, location, errorType)
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Record success metrics
	if h.metrics != nil {
		h.metrics.RecordRelease(resp.SKU, resp.Location, "success")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// AdjustStock handles POST /adjust
func (h *ReservationHandlers) AdjustStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SKU       string `json:"sku"`
		Location  string `json:"location"`
		Quantity  int64  `json:"quantity"`
		Reason    string `json:"reason"`
		CreatedBy string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		if h.metrics != nil {
			h.metrics.RecordAdjustmentError(req.SKU, req.Location, "invalid_request")
		}
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	appReq := app.AdjustStockRequest{
		SKU:       req.SKU,
		Location:  req.Location,
		Quantity:  req.Quantity,
		Reason:    req.Reason,
		CreatedBy: req.CreatedBy,
	}

	resp, err := h.adjustUseCase.Execute(r.Context(), appReq)

	if err != nil {
		h.logger.Error("failed to adjust stock", zap.Error(err))

		// Record error metrics
		if h.metrics != nil {
			errorType := "unknown"
			if strings.Contains(err.Error(), "negative") {
				errorType = "negative_stock"
			}
			h.metrics.RecordAdjustmentError(req.SKU, req.Location, errorType)
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Record success metrics
	if h.metrics != nil {
		adjustmentType := "increase"
		if req.Quantity < 0 {
			adjustmentType = "decrease"
		}
		h.metrics.RecordAdjustment(resp.SKU, resp.Location, adjustmentType)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// QueryAvailable handles GET /available
func (h *ReservationHandlers) QueryAvailable(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sku := r.URL.Query().Get("sku")
	location := r.URL.Query().Get("location")

	if sku == "" || location == "" {
		http.Error(w, "sku and location query parameters are required", http.StatusBadRequest)
		return
	}

	appReq := app.QueryAvailableRequest{
		SKU:      sku,
		Location: location,
	}

	resp, err := h.queryUseCase.Execute(r.Context(), appReq)
	duration := time.Since(startTime).Seconds()

	if err != nil {
		h.logger.Error("failed to query available stock", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Record query metrics
	// Note: cache hit detection would need to be added to QueryAvailableUseCase
	if h.metrics != nil {
		cacheHit := false // TODO: Get from use case
		h.metrics.RecordQuery(sku, location, cacheHit, duration)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

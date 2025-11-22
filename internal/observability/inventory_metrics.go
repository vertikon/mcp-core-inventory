// Package observability provides inventory-specific Prometheus metrics
package observability

import (
	"github.com/prometheus/client_golang/prometheus"
)

// InventoryMetrics provides inventory-specific metrics
type InventoryMetrics struct {
	// Stock reservation metrics
	stockReservationTotal    *prometheus.CounterVec
	stockReservationErrors   *prometheus.CounterVec
	stockReservationDuration *prometheus.HistogramVec

	// Stock confirmation metrics
	stockConfirmationTotal    *prometheus.CounterVec
	stockConfirmationErrors   *prometheus.CounterVec
	stockConfirmationDuration *prometheus.HistogramVec

	// Stock release metrics
	stockReleaseTotal  *prometheus.CounterVec
	stockReleaseErrors *prometheus.CounterVec

	// Stock adjustment metrics
	stockAdjustmentTotal  *prometheus.CounterVec
	stockAdjustmentErrors *prometheus.CounterVec

	// Stock query metrics
	stockQueryTotal    *prometheus.CounterVec
	stockQueryDuration *prometheus.HistogramVec

	// Race condition incidents (CRITICAL)
	raceConditionIncidents prometheus.Counter

	// Redis metrics
	redisLatencyMs             *prometheus.HistogramVec
	redisLockAcquisitionTotal  *prometheus.CounterVec
	redisLockAcquisitionErrors *prometheus.CounterVec

	// Postgres metrics
	postgresLockWaitMs    *prometheus.HistogramVec
	postgresQueryDuration *prometheus.HistogramVec

	// Reservation TTL metrics
	reservationsTTLExpiredTotal prometheus.Counter
	reservationsActive          prometheus.Gauge

	// Idempotency metrics
	idempotencyHitsTotal   prometheus.Counter
	idempotencyMissesTotal prometheus.Counter
}

// NewInventoryMetrics creates new inventory-specific metrics
func NewInventoryMetrics() *InventoryMetrics {
	m := &InventoryMetrics{
		// Stock reservation metrics
		stockReservationTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "stock_reservation_total",
				Help: "Total number of stock reservations",
			},
			[]string{"sku", "location", "status"},
		),
		stockReservationErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "stock_reservation_errors_total",
				Help: "Total number of stock reservation errors",
			},
			[]string{"sku", "location", "error_type"},
		),
		stockReservationDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "stock_reservation_duration_seconds",
				Help:    "Stock reservation duration in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
			},
			[]string{"sku", "location"},
		),

		// Stock confirmation metrics
		stockConfirmationTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "stock_confirmation_total",
				Help: "Total number of stock confirmations",
			},
			[]string{"sku", "location", "status"},
		),
		stockConfirmationErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "stock_confirmation_errors_total",
				Help: "Total number of stock confirmation errors",
			},
			[]string{"sku", "location", "error_type"},
		),
		stockConfirmationDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "stock_confirmation_duration_seconds",
				Help:    "Stock confirmation duration in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
			},
			[]string{"sku", "location"},
		),

		// Stock release metrics
		stockReleaseTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "stock_release_total",
				Help: "Total number of stock releases",
			},
			[]string{"sku", "location", "status"},
		),
		stockReleaseErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "stock_release_errors_total",
				Help: "Total number of stock release errors",
			},
			[]string{"sku", "location", "error_type"},
		),

		// Stock adjustment metrics
		stockAdjustmentTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "stock_adjustment_total",
				Help: "Total number of stock adjustments",
			},
			[]string{"sku", "location", "type"},
		),
		stockAdjustmentErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "stock_adjustment_errors_total",
				Help: "Total number of stock adjustment errors",
			},
			[]string{"sku", "location", "error_type"},
		),

		// Stock query metrics
		stockQueryTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "stock_query_total",
				Help: "Total number of stock queries",
			},
			[]string{"sku", "location", "cache_hit"},
		),
		stockQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "stock_query_duration_seconds",
				Help:    "Stock query duration in seconds",
				Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5},
			},
			[]string{"sku", "location", "cache_hit"},
		),

		// Race condition incidents (CRITICAL - P0 alert)
		raceConditionIncidents: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "race_condition_incidents_total",
				Help: "Total number of race condition incidents detected (P0 alert if > 0)",
			},
		),

		// Redis metrics
		redisLatencyMs: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "redis_latency_ms",
				Help:    "Redis operation latency in milliseconds",
				Buckets: []float64{0.1, 0.5, 1.0, 2.5, 5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0},
			},
			[]string{"operation"},
		),
		redisLockAcquisitionTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "redis_lock_acquisition_total",
				Help: "Total number of Redis lock acquisition attempts",
			},
			[]string{"status"},
		),
		redisLockAcquisitionErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "redis_lock_acquisition_errors_total",
				Help: "Total number of Redis lock acquisition errors",
			},
			[]string{"error_type"},
		),

		// Postgres metrics
		postgresLockWaitMs: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "postgres_lock_wait_ms",
				Help:    "PostgreSQL lock wait time in milliseconds",
				Buckets: []float64{0.1, 0.5, 1.0, 2.5, 5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0, 1000.0, 2000.0},
			},
			[]string{"table", "operation"},
		),
		postgresQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "postgres_query_duration_seconds",
				Help:    "PostgreSQL query duration in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
			},
			[]string{"operation"},
		),

		// Reservation TTL metrics
		reservationsTTLExpiredTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "reservations_ttl_expired_total",
				Help: "Total number of reservations expired by TTL",
			},
		),
		reservationsActive: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "reservations_active",
				Help: "Current number of active reservations",
			},
		),

		// Idempotency metrics
		idempotencyHitsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "idempotency_hits_total",
				Help: "Total number of idempotency key hits (duplicate requests)",
			},
		),
		idempotencyMissesTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "idempotency_misses_total",
				Help: "Total number of idempotency key misses (new requests)",
			},
		),
	}

	// Register all metrics
	prometheus.MustRegister(
		m.stockReservationTotal,
		m.stockReservationErrors,
		m.stockReservationDuration,
		m.stockConfirmationTotal,
		m.stockConfirmationErrors,
		m.stockConfirmationDuration,
		m.stockReleaseTotal,
		m.stockReleaseErrors,
		m.stockAdjustmentTotal,
		m.stockAdjustmentErrors,
		m.stockQueryTotal,
		m.stockQueryDuration,
		m.raceConditionIncidents,
		m.redisLatencyMs,
		m.redisLockAcquisitionTotal,
		m.redisLockAcquisitionErrors,
		m.postgresLockWaitMs,
		m.postgresQueryDuration,
		m.reservationsTTLExpiredTotal,
		m.reservationsActive,
		m.idempotencyHitsTotal,
		m.idempotencyMissesTotal,
	)

	return m
}

// RecordReservation records a stock reservation
func (m *InventoryMetrics) RecordReservation(sku, location, status string, duration float64) {
	m.stockReservationTotal.WithLabelValues(sku, location, status).Inc()
	m.stockReservationDuration.WithLabelValues(sku, location).Observe(duration)
}

// RecordReservationError records a stock reservation error
func (m *InventoryMetrics) RecordReservationError(sku, location, errorType string) {
	m.stockReservationErrors.WithLabelValues(sku, location, errorType).Inc()
}

// RecordConfirmation records a stock confirmation
func (m *InventoryMetrics) RecordConfirmation(sku, location, status string, duration float64) {
	m.stockConfirmationTotal.WithLabelValues(sku, location, status).Inc()
	m.stockConfirmationDuration.WithLabelValues(sku, location).Observe(duration)
}

// RecordConfirmationError records a stock confirmation error
func (m *InventoryMetrics) RecordConfirmationError(sku, location, errorType string) {
	m.stockConfirmationErrors.WithLabelValues(sku, location, errorType).Inc()
}

// RecordRelease records a stock release
func (m *InventoryMetrics) RecordRelease(sku, location, status string) {
	m.stockReleaseTotal.WithLabelValues(sku, location, status).Inc()
}

// RecordReleaseError records a stock release error
func (m *InventoryMetrics) RecordReleaseError(sku, location, errorType string) {
	m.stockReleaseErrors.WithLabelValues(sku, location, errorType).Inc()
}

// RecordAdjustment records a stock adjustment
func (m *InventoryMetrics) RecordAdjustment(sku, location, adjustmentType string) {
	m.stockAdjustmentTotal.WithLabelValues(sku, location, adjustmentType).Inc()
}

// RecordAdjustmentError records a stock adjustment error
func (m *InventoryMetrics) RecordAdjustmentError(sku, location, errorType string) {
	m.stockAdjustmentErrors.WithLabelValues(sku, location, errorType).Inc()
}

// RecordQuery records a stock query
func (m *InventoryMetrics) RecordQuery(sku, location string, cacheHit bool, duration float64) {
	cacheHitStr := "false"
	if cacheHit {
		cacheHitStr = "true"
	}
	m.stockQueryTotal.WithLabelValues(sku, location, cacheHitStr).Inc()
	m.stockQueryDuration.WithLabelValues(sku, location, cacheHitStr).Observe(duration)
}

// RecordRaceCondition records a race condition incident (P0 alert)
func (m *InventoryMetrics) RecordRaceCondition() {
	m.raceConditionIncidents.Inc()
}

// RecordRedisLatency records Redis operation latency
func (m *InventoryMetrics) RecordRedisLatency(operation string, latencyMs float64) {
	m.redisLatencyMs.WithLabelValues(operation).Observe(latencyMs)
}

// RecordRedisLockAcquisition records Redis lock acquisition
func (m *InventoryMetrics) RecordRedisLockAcquisition(status string) {
	m.redisLockAcquisitionTotal.WithLabelValues(status).Inc()
}

// RecordRedisLockError records Redis lock acquisition error
func (m *InventoryMetrics) RecordRedisLockError(errorType string) {
	m.redisLockAcquisitionErrors.WithLabelValues(errorType).Inc()
}

// RecordPostgresLockWait records PostgreSQL lock wait time
func (m *InventoryMetrics) RecordPostgresLockWait(table, operation string, waitMs float64) {
	m.postgresLockWaitMs.WithLabelValues(table, operation).Observe(waitMs)
}

// RecordPostgresQueryDuration records PostgreSQL query duration
func (m *InventoryMetrics) RecordPostgresQueryDuration(operation string, duration float64) {
	m.postgresQueryDuration.WithLabelValues(operation).Observe(duration)
}

// RecordReservationTTLExpired records an expired reservation
func (m *InventoryMetrics) RecordReservationTTLExpired() {
	m.reservationsTTLExpiredTotal.Inc()
	m.reservationsActive.Dec()
}

// SetReservationsActive sets the number of active reservations
func (m *InventoryMetrics) SetReservationsActive(count float64) {
	m.reservationsActive.Set(count)
}

// RecordIdempotencyHit records an idempotency key hit
func (m *InventoryMetrics) RecordIdempotencyHit() {
	m.idempotencyHitsTotal.Inc()
}

// RecordIdempotencyMiss records an idempotency key miss
func (m *InventoryMetrics) RecordIdempotencyMiss() {
	m.idempotencyMissesTotal.Inc()
}

package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/vertikon/mcp-core-inventory/internal/app"
	"github.com/vertikon/mcp-core-inventory/internal/domain/ledger"
	"go.uber.org/zap"
)

// EventPublisher implements event publishing using NATS JetStream
type EventPublisher struct {
	js     nats.JetStreamContext
	logger *zap.Logger
}

// NewEventPublisher creates a new NATS event publisher
func NewEventPublisher(js nats.JetStreamContext, logger *zap.Logger) *EventPublisher {
	return &EventPublisher{
		js:     js,
		logger: logger,
	}
}

// PublishReserveRequest publishes an inventory.reserve.request event
func (p *EventPublisher) PublishReserveRequest(ctx context.Context, req app.ReserveStockRequest, reservationID string) error {
	event := ReserveRequestEvent{
		ReservationID:  reservationID,
		IdempotencyKey: req.IdempotencyKey,
		SKU:            req.SKU,
		Location:       req.Location,
		Quantity:       req.Quantity,
		Timestamp:      fmt.Sprintf("%d", ctx.Value("timestamp").(int64)),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = p.js.Publish("inventory.reserve.request", data)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.Info("published inventory.reserve.request",
		zap.String("reservation_id", reservationID),
		zap.String("sku", req.SKU),
		zap.String("location", req.Location),
	)

	return nil
}

// PublishReserveConfirmed publishes an inventory.reserve.confirmed event
func (p *EventPublisher) PublishReserveConfirmed(ctx context.Context, reservation *ledger.Reservation) error {
	event := ReserveConfirmedEvent{
		ReservationID: reservation.ID(),
		SKU:           reservation.SKU(),
		Location:      reservation.Location(),
		Quantity:      reservation.Quantity(),
		ConfirmedAt:   reservation.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
		Timestamp:     fmt.Sprintf("%d", ctx.Value("timestamp").(int64)),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = p.js.Publish("inventory.reserve.confirmed", data)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.Info("published inventory.reserve.confirmed",
		zap.String("reservation_id", reservation.ID()),
		zap.String("sku", reservation.SKU()),
	)

	return nil
}

// ReserveRequestEvent represents the inventory.reserve.request event
type ReserveRequestEvent struct {
	ReservationID  string `json:"reservation_id"`
	IdempotencyKey string `json:"idempotency_key"`
	SKU            string `json:"sku"`
	Location       string `json:"location"`
	Quantity       int64  `json:"quantity"`
	Timestamp      string `json:"timestamp"`
}

// ReserveConfirmedEvent represents the inventory.reserve.confirmed event
type ReserveConfirmedEvent struct {
	ReservationID string `json:"reservation_id"`
	SKU           string `json:"sku"`
	Location      string `json:"location"`
	Quantity      int64  `json:"quantity"`
	ConfirmedAt   string `json:"confirmed_at"`
	Timestamp     string `json:"timestamp"`
}

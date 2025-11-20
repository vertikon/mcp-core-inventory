package batches

import (
	"errors"
	"time"
)

// Batch represents a batch/lot with expiration date
type Batch struct {
	id          string
	sku         string
	location    string
	quantity    int64
	batchNumber string
	expiresAt   *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

// NewBatch creates a new batch
func NewBatch(id, sku, location, batchNumber string, quantity int64, expiresAt *time.Time) (*Batch, error) {
	if id == "" {
		return nil, errors.New("batch id cannot be empty")
	}
	if sku == "" {
		return nil, errors.New("sku cannot be empty")
	}
	if location == "" {
		return nil, errors.New("location cannot be empty")
	}
	if batchNumber == "" {
		return nil, errors.New("batch number cannot be empty")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	now := time.Now()
	return &Batch{
		id:          id,
		sku:         sku,
		location:    location,
		quantity:    quantity,
		batchNumber: batchNumber,
		expiresAt:   expiresAt,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ID returns the batch ID
func (b *Batch) ID() string {
	return b.id
}

// SKU returns the SKU
func (b *Batch) SKU() string {
	return b.sku
}

// Location returns the location
func (b *Batch) Location() string {
	return b.location
}

// Quantity returns the batch quantity
func (b *Batch) Quantity() int64 {
	return b.quantity
}

// BatchNumber returns the batch number
func (b *Batch) BatchNumber() string {
	return b.batchNumber
}

// ExpiresAt returns the expiration date
func (b *Batch) ExpiresAt() *time.Time {
	return b.expiresAt
}

// CreatedAt returns the creation timestamp
func (b *Batch) CreatedAt() time.Time {
	return b.createdAt
}

// UpdatedAt returns the last update timestamp
func (b *Batch) UpdatedAt() time.Time {
	return b.updatedAt
}

// IsExpired checks if the batch has expired
func (b *Batch) IsExpired() bool {
	if b.expiresAt == nil {
		return false
	}
	return time.Now().After(*b.expiresAt)
}

// Consume reduces the batch quantity
func (b *Batch) Consume(quantity int64) error {
	if quantity <= 0 {
		return errors.New("consumption quantity must be positive")
	}
	if b.quantity < quantity {
		return errors.New("insufficient batch quantity")
	}

	b.quantity -= quantity
	b.updatedAt = time.Now()
	return nil
}


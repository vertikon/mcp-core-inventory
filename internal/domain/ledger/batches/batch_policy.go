package batches

import (
	"errors"
	"sort"
)

// ConsumptionPolicy defines how batches are consumed
type ConsumptionPolicy string

const (
	// FIFO - First In First Out
	FIFO ConsumptionPolicy = "fifo"
	// FEFO - First Expiry First Out
	FEFO ConsumptionPolicy = "fefo"
)

// BatchPolicy enforces batch consumption policies
type BatchPolicy struct {
	policy ConsumptionPolicy
}

// NewBatchPolicy creates a new batch policy
func NewBatchPolicy(policy ConsumptionPolicy) (*BatchPolicy, error) {
	if policy != FIFO && policy != FEFO {
		return nil, errors.New("invalid consumption policy")
	}
	return &BatchPolicy{policy: policy}, nil
}

// SelectBatchesForConsumption selects batches to consume based on the policy
// Returns batches and quantities to consume
func (bp *BatchPolicy) SelectBatchesForConsumption(batches []*Batch, quantity int64) ([]*Batch, []int64, error) {
	if quantity <= 0 {
		return nil, nil, errors.New("quantity must be positive")
	}

	// Filter available batches (not expired, quantity > 0)
	available := make([]*Batch, 0)
	for _, b := range batches {
		if !b.IsExpired() && b.Quantity() > 0 {
			available = append(available, b)
		}
	}

	if len(available) == 0 {
		return nil, nil, errors.New("no available batches")
	}

	// Sort based on policy
	switch bp.policy {
	case FIFO:
		sort.Slice(available, func(i, j int) bool {
			return available[i].CreatedAt().Before(available[j].CreatedAt())
		})
	case FEFO:
		sort.Slice(available, func(i, j int) bool {
			expI := available[i].ExpiresAt()
			expJ := available[j].ExpiresAt()
			if expI == nil && expJ == nil {
				return available[i].CreatedAt().Before(available[j].CreatedAt())
			}
			if expI == nil {
				return false
			}
			if expJ == nil {
				return true
			}
			return expI.Before(*expJ)
		})
	}

	// Select batches to consume
	selectedBatches := make([]*Batch, 0)
	quantities := make([]int64, 0)
	remaining := quantity

	for _, batch := range available {
		if remaining <= 0 {
			break
		}
		consumeQty := batch.Quantity()
		if consumeQty > remaining {
			consumeQty = remaining
		}
		selectedBatches = append(selectedBatches, batch)
		quantities = append(quantities, consumeQty)
		remaining -= consumeQty
	}

	if remaining > 0 {
		return nil, nil, errors.New("insufficient batch quantity")
	}

	return selectedBatches, quantities, nil
}

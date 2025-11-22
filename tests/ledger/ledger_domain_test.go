package ledger_test

import (
	"testing"
	"time"

	"github.com/vertikon/mcp-core-inventory/internal/domain/ledger"
)

func TestStockLedger_Reserve(t *testing.T) {
	tests := []struct {
		name          string
		initialQty    int64
		reserveQty    int64
		wantErr       bool
		wantAvailable int64
		wantReserved  int64
	}{
		{
			name:          "successful reservation",
			initialQty:    100,
			reserveQty:    50,
			wantErr:       false,
			wantAvailable: 50,
			wantReserved:  50,
		},
		{
			name:          "insufficient stock",
			initialQty:    100,
			reserveQty:    150,
			wantErr:       true,
			wantAvailable: 100,
			wantReserved:  0,
		},
		{
			name:          "zero quantity",
			initialQty:    100,
			reserveQty:    0,
			wantErr:       true,
			wantAvailable: 100,
			wantReserved:  0,
		},
		{
			name:          "negative quantity",
			initialQty:    100,
			reserveQty:    -10,
			wantErr:       true,
			wantAvailable: 100,
			wantReserved:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ledger, err := ledger.NewStockLedger("test-id", "SKU-001", "LOC-001", tt.initialQty)
			if err != nil {
				t.Fatalf("Failed to create ledger: %v", err)
			}

			err = ledger.Reserve(tt.reserveQty)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reserve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if ledger.Available() != tt.wantAvailable {
					t.Errorf("Available() = %v, want %v", ledger.Available(), tt.wantAvailable)
				}
				if ledger.Reserved() != tt.wantReserved {
					t.Errorf("Reserved() = %v, want %v", ledger.Reserved(), tt.wantReserved)
				}
			}
		})
	}
}

func TestStockLedger_ConfirmReservation(t *testing.T) {
	tests := []struct {
		name         string
		initialQty   int64
		reserveQty   int64
		confirmQty   int64
		wantErr      bool
		wantTotal    int64
		wantReserved int64
	}{
		{
			name:         "successful confirmation",
			initialQty:   100,
			reserveQty:   50,
			confirmQty:   50,
			wantErr:      false,
			wantTotal:    50,
			wantReserved: 0,
		},
		{
			name:         "insufficient reserved",
			initialQty:   100,
			reserveQty:   50,
			confirmQty:   60,
			wantErr:      true,
			wantTotal:    100,
			wantReserved: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ledger, _ := ledger.NewStockLedger("test-id", "SKU-001", "LOC-001", tt.initialQty)
			ledger.Reserve(tt.reserveQty)

			err := ledger.ConfirmReservation(tt.confirmQty)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfirmReservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if ledger.Total() != tt.wantTotal {
					t.Errorf("Total() = %v, want %v", ledger.Total(), tt.wantTotal)
				}
				if ledger.Reserved() != tt.wantReserved {
					t.Errorf("Reserved() = %v, want %v", ledger.Reserved(), tt.wantReserved)
				}
			}
		})
	}
}

func TestReservation_IsExpired(t *testing.T) {
	tests := []struct {
		name        string
		ttl         time.Duration
		waitTime    time.Duration
		wantExpired bool
	}{
		{
			name:        "not expired",
			ttl:         5 * time.Minute,
			waitTime:    0,
			wantExpired: false,
		},
		{
			name:        "expired",
			ttl:         100 * time.Millisecond,
			waitTime:    150 * time.Millisecond,
			wantExpired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reservation, err := ledger.NewReservation(
				"res-001",
				"ledger-001",
				"SKU-001",
				"LOC-001",
				"idempotency-key",
				10,
				tt.ttl,
			)
			if err != nil {
				t.Fatalf("Failed to create reservation: %v", err)
			}

			if tt.waitTime > 0 {
				time.Sleep(tt.waitTime)
			}

			if reservation.IsExpired() != tt.wantExpired {
				t.Errorf("IsExpired() = %v, want %v", reservation.IsExpired(), tt.wantExpired)
			}
		})
	}
}

func TestPolicy_ValidateReservation(t *testing.T) {
	policy := ledger.NewPolicy()

	tests := []struct {
		name      string
		available int64
		quantity  int64
		wantErr   bool
	}{
		{
			name:      "valid reservation",
			available: 100,
			quantity:  50,
			wantErr:   false,
		},
		{
			name:      "insufficient stock",
			available: 100,
			quantity:  150,
			wantErr:   true,
		},
		{
			name:      "zero quantity",
			available: 100,
			quantity:  0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ledger, _ := ledger.NewStockLedger("test-id", "SKU-001", "LOC-001", tt.available)
			err := policy.ValidateReservation(ledger, tt.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateReservation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

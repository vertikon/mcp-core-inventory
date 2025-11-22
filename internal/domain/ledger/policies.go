package ledger

import "errors"

// Policy enforces business rules for inventory operations
type Policy struct{}

// NewPolicy creates a new policy instance
func NewPolicy() *Policy {
	return &Policy{}
}

// ValidateReservation validates if a reservation can be made
// Enforces: "Não vender sem saldo"
func (p *Policy) ValidateReservation(ledger *StockLedger, quantity int64) error {
	if ledger == nil {
		return errors.New("ledger cannot be nil")
	}
	if quantity <= 0 {
		return errors.New("reservation quantity must be positive")
	}
	if !ledger.CanReserve(quantity) {
		return ErrInsufficientStock
	}
	return nil
}

// ValidateConfirmation validates if a reservation can be confirmed
// Enforces: ordem correta reserva → confirmação/expedição → baixa
func (p *Policy) ValidateConfirmation(reservation *Reservation) error {
	if reservation == nil {
		return errors.New("reservation cannot be nil")
	}
	if reservation.Status() != ReservationStatusPending {
		return errors.New("only pending reservations can be confirmed")
	}
	if reservation.IsExpired() {
		return ErrReservationExpired
	}
	return nil
}

// ValidateAdjustment validates if an adjustment can be made
// Enforces: saldo nunca negativo
func (p *Policy) ValidateAdjustment(ledger *StockLedger, quantity int64) error {
	if ledger == nil {
		return errors.New("ledger cannot be nil")
	}
	newTotal := ledger.Total() + quantity
	if newTotal < 0 {
		return ErrNegativeStock
	}
	return nil
}

// ValidateRelease validates if a reservation can be released
func (p *Policy) ValidateRelease(reservation *Reservation) error {
	if reservation == nil {
		return errors.New("reservation cannot be nil")
	}
	if reservation.Status() != ReservationStatusPending {
		return errors.New("only pending reservations can be released")
	}
	return nil
}

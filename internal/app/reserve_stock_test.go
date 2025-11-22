package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vertikon/mcp-core-inventory/internal/domain/ledger"
)

type testLedgerRepo struct {
	ledger       *ledger.StockLedger
	reservations map[string]*ledger.Reservation
}

func (r *testLedgerRepo) FindBySKUAndLocation(ctx context.Context, sku, location string) (*ledger.StockLedger, error) {
	if r.ledger == nil {
		return nil, errors.New("not found")
	}
	return r.ledger, nil
}

func (r *testLedgerRepo) Save(ctx context.Context, l *ledger.StockLedger) error {
	r.ledger = l
	return nil
}

func (r *testLedgerRepo) SaveReservation(ctx context.Context, reservation *ledger.Reservation) error {
	if r.reservations == nil {
		r.reservations = make(map[string]*ledger.Reservation)
	}
	r.reservations[reservation.IdempotencyKey()] = reservation
	return nil
}

func (r *testLedgerRepo) FindReservationByKey(ctx context.Context, idempotencyKey string) (*ledger.Reservation, error) {
	if res, ok := r.reservations[idempotencyKey]; ok {
		return res, nil
	}
	return nil, errors.New("not found")
}

type testLock struct {
	acquired bool
}

func (l *testLock) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if l.acquired {
		return false, nil
	}
	l.acquired = true
	return true, nil
}

func (l *testLock) ReleaseLock(ctx context.Context, key string) error {
	l.acquired = false
	return nil
}

type testPublisher struct {
	reserveEventCalled bool
}

func (p *testPublisher) PublishReserveRequest(ctx context.Context, req ReserveStockRequest, reservationID string) error {
	p.reserveEventCalled = true
	return nil
}

func (p *testPublisher) PublishReserveConfirmed(ctx context.Context, reservation *ledger.Reservation) error {
	return nil
}

func TestReserveStockUseCase_Success(t *testing.T) {
	stock, err := ledger.NewStockLedger("ledger-1", "SKU1", "CD-SP", 10)
	if err != nil {
		t.Fatalf("failed to create ledger: %v", err)
	}
	repo := &testLedgerRepo{ledger: stock}
	lock := &testLock{}
	pub := &testPublisher{}
	uc := NewReserveStockUseCase(repo, lock, pub, nil)

	resp, err := uc.Execute(context.Background(), ReserveStockRequest{
		SKU:      "SKU1",
		Location: "CD-SP",
		Quantity: 5,
	})
	if err != nil {
		t.Fatalf("expected success, got error %v", err)
	}
	if resp.Quantity != 5 {
		t.Fatalf("expected quantity 5, got %d", resp.Quantity)
	}
	if !pub.reserveEventCalled {
		t.Fatalf("expected publisher to be called")
	}
	if repo.ledger.Reserved() != 5 {
		t.Fatalf("expected 5 reserved, got %d", repo.ledger.Reserved())
	}
}

func TestReserveStockUseCase_Idempotency(t *testing.T) {
	stock, _ := ledger.NewStockLedger("ledger-1", "SKU1", "CD-SP", 10)
	repo := &testLedgerRepo{ledger: stock, reservations: make(map[string]*ledger.Reservation)}
	lock := &testLock{}
	pub := &testPublisher{}
	uc := NewReserveStockUseCase(repo, lock, pub, nil)

	ctx := context.Background()
	req := ReserveStockRequest{
		SKU:            "SKU1",
		Location:       "CD-SP",
		Quantity:       2,
		IdempotencyKey: "token-1",
	}

	firstResp, err := uc.Execute(ctx, req)
	if err != nil {
		t.Fatalf("expected first reservation to succeed: %v", err)
	}

	secondResp, err := uc.Execute(ctx, req)
	if err != nil {
		t.Fatalf("expected idempotent reservation to succeed: %v", err)
	}

	if firstResp.ReservationID != secondResp.ReservationID {
		t.Fatalf("expected same reservation id for idempotent call")
	}
}

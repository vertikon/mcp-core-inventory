-- Migration 0001: Create stock ledger and reservations tables
-- This migration creates the core tables for inventory management

-- Stock Ledger table
-- Maintains the ACID ledger of stock balances
CREATE TABLE IF NOT EXISTS stock_ledger (
    id VARCHAR(255) PRIMARY KEY,
    sku VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    available BIGINT NOT NULL DEFAULT 0,
    reserved BIGINT NOT NULL DEFAULT 0,
    total BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(sku, location)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_stock_ledger_sku_location ON stock_ledger(sku, location);
CREATE INDEX IF NOT EXISTS idx_stock_ledger_sku ON stock_ledger(sku);
CREATE INDEX IF NOT EXISTS idx_stock_ledger_location ON stock_ledger(location);

-- Reservations table
-- Tracks stock reservations with TTL and idempotency
CREATE TABLE IF NOT EXISTS reservations (
    id VARCHAR(255) PRIMARY KEY,
    ledger_id VARCHAR(255) NOT NULL,
    sku VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    quantity BIGINT NOT NULL,
    idempotency_key VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ledger_id) REFERENCES stock_ledger(id) ON DELETE CASCADE
);

-- Indexes for reservations
CREATE INDEX IF NOT EXISTS idx_reservations_idempotency_key ON reservations(idempotency_key);
CREATE INDEX IF NOT EXISTS idx_reservations_ledger_id ON reservations(ledger_id);
CREATE INDEX IF NOT EXISTS idx_reservations_sku_location ON reservations(sku, location);
CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations(status);
CREATE INDEX IF NOT EXISTS idx_reservations_expires_at ON reservations(expires_at);

-- Stock Movements table
-- Audit trail of all stock movements
CREATE TABLE IF NOT EXISTS stock_movements (
    id VARCHAR(255) PRIMARY KEY,
    ledger_id VARCHAR(255) NOT NULL,
    movement_type VARCHAR(50) NOT NULL,
    quantity BIGINT NOT NULL,
    reference_id VARCHAR(255),
    reason TEXT,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    FOREIGN KEY (ledger_id) REFERENCES stock_ledger(id) ON DELETE CASCADE
);

-- Indexes for movements
CREATE INDEX IF NOT EXISTS idx_stock_movements_ledger_id ON stock_movements(ledger_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_type ON stock_movements(movement_type);
CREATE INDEX IF NOT EXISTS idx_stock_movements_reference_id ON stock_movements(reference_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_created_at ON stock_movements(created_at);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers to auto-update updated_at
CREATE TRIGGER update_stock_ledger_updated_at BEFORE UPDATE ON stock_ledger
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reservations_updated_at BEFORE UPDATE ON reservations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


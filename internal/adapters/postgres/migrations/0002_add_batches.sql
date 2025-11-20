-- Migration 0002: Add batch/lot tracking with expiration dates
-- This migration adds support for FIFO/FEFO batch consumption policies

-- Batches table
-- Tracks batches/lots with expiration dates for FIFO/FEFO policies
CREATE TABLE IF NOT EXISTS batches (
    id VARCHAR(255) PRIMARY KEY,
    sku VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    quantity BIGINT NOT NULL,
    batch_number VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(sku, location, batch_number)
);

-- Indexes for batches
CREATE INDEX IF NOT EXISTS idx_batches_sku_location ON batches(sku, location);
CREATE INDEX IF NOT EXISTS idx_batches_batch_number ON batches(batch_number);
CREATE INDEX IF NOT EXISTS idx_batches_expires_at ON batches(expires_at);
CREATE INDEX IF NOT EXISTS idx_batches_sku_location_expires ON batches(sku, location, expires_at);

-- Batch Movements table
-- Tracks batch consumption for audit trail
CREATE TABLE IF NOT EXISTS batch_movements (
    id VARCHAR(255) PRIMARY KEY,
    batch_id VARCHAR(255) NOT NULL,
    movement_type VARCHAR(50) NOT NULL,
    quantity BIGINT NOT NULL,
    reference_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    FOREIGN KEY (batch_id) REFERENCES batches(id) ON DELETE CASCADE
);

-- Indexes for batch movements
CREATE INDEX IF NOT EXISTS idx_batch_movements_batch_id ON batch_movements(batch_id);
CREATE INDEX IF NOT EXISTS idx_batch_movements_type ON batch_movements(movement_type);

-- Trigger to auto-update updated_at for batches
CREATE TRIGGER update_batches_updated_at BEFORE UPDATE ON batches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- View for available batches (not expired, quantity > 0)
CREATE OR REPLACE VIEW available_batches AS
SELECT 
    id,
    sku,
    location,
    quantity,
    batch_number,
    expires_at,
    created_at,
    updated_at
FROM batches
WHERE quantity > 0 
    AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP);


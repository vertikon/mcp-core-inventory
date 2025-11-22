-- Script de Inicialização do Banco de Dados - mcp-core-inventory
-- Este script é executado automaticamente na criação do container PostgreSQL

-- Extensões necessárias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Schema principal
CREATE SCHEMA IF NOT EXISTS inventory;
CREATE SCHEMA IF NOT EXISTS events;
CREATE SCHEMA IF NOT EXISTS audit;

-- Configuração de permissões
GRANT USAGE ON SCHEMA inventory TO postgres;
GRANT USAGE ON SCHEMA events TO postgres;
GRANT USAGE ON SCHEMA audit TO postgres;

-- Tabelas de Inventário (Core Ledger)
CREATE TABLE IF NOT EXISTS inventory.products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sku VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(500) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    unit_cost DECIMAL(15,2),
    unit_price DECIMAL(15,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version INTEGER DEFAULT 1
);

CREATE TABLE IF NOT EXISTS inventory.warehouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    location JSONB,
    capacity INTEGER,
    manager_id UUID,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS inventory.inventory_ledger (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES inventory.products(id),
    warehouse_id UUID NOT NULL REFERENCES inventory.warehouses(id),
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('IN', 'OUT', 'ADJUSTMENT', 'TRANSFER')),
    quantity INTEGER NOT NULL,
    unit_cost DECIMAL(15,2),
    reference_id UUID,
    reference_type VARCHAR(50),
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS inventory.inventory_balances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES inventory.products(id),
    warehouse_id UUID NOT NULL REFERENCES inventory.warehouses(id),
    quantity_on_hand INTEGER NOT NULL DEFAULT 0,
    quantity_reserved INTEGER NOT NULL DEFAULT 0,
    quantity_available INTEGER GENERATED ALWAYS AS (quantity_on_hand - quantity_reserved) STORED,
    unit_cost DECIMAL(15,2),
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(product_id, warehouse_id)
);

-- Tabelas de Eventos
CREATE TABLE IF NOT EXISTS events.events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    event_data JSONB NOT NULL,
    event_version INTEGER DEFAULT 1,
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB
);

CREATE TABLE IF NOT EXISTS events.event_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_name VARCHAR(100) NOT NULL,
    event_types TEXT[] NOT NULL,
    endpoint_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabelas de Auditoria
CREATE TABLE IF NOT EXISTS audit.audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_inventory_ledger_product_warehouse ON inventory.inventory_ledger(product_id, warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventory_ledger_transaction_date ON inventory.inventory_ledger(transaction_date);
CREATE INDEX IF NOT EXISTS idx_inventory_balances_product ON inventory.inventory_balances(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_balances_warehouse ON inventory.inventory_balances(warehouse_id);

CREATE INDEX IF NOT EXISTS idx_events_aggregate ON events.events(aggregate_id, aggregate_type);
CREATE INDEX IF NOT EXISTS idx_events_type_occurred ON events.events(event_type, occurred_at);
CREATE INDEX IF NOT EXISTS idx_events_processed ON events.events(processed_at) WHERE processed_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_audit_user_timestamp ON audit.audit_logs(user_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_audit_resource ON audit.audit_logs(resource_type, resource_id);

-- Triggers para atualização de timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON inventory.products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_warehouses_updated_at BEFORE UPDATE ON inventory.warehouses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Trigger para manter saldos de inventário
CREATE OR REPLACE FUNCTION update_inventory_balance()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO inventory.inventory_balances (product_id, warehouse_id, quantity_on_hand, unit_cost)
        VALUES (NEW.product_id, NEW.warehouse_id, 
                CASE WHEN NEW.transaction_type = 'IN' THEN NEW.quantity 
                     WHEN NEW.transaction_type = 'OUT' THEN -NEW.quantity 
                     ELSE NEW.quantity END,
                NEW.unit_cost)
        ON CONFLICT (product_id, warehouse_id) DO UPDATE SET
            quantity_on_hand = inventory.inventory_balances.quantity_on_hand + 
                CASE WHEN NEW.transaction_type = 'IN' THEN NEW.quantity 
                     WHEN NEW.transaction_type = 'OUT' THEN -NEW.quantity 
                     ELSE NEW.quantity END,
            unit_cost = NEW.unit_cost,
            last_updated = NOW();
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_inventory_balance
    AFTER INSERT ON inventory.inventory_ledger
    FOR EACH ROW EXECUTE FUNCTION update_inventory_balance();

-- Dados iniciais (opcional)
INSERT INTO inventory.warehouses (code, name, location, capacity) VALUES
('WH-001', 'Principal Warehouse', '{"city": "São Paulo", "state": "SP"}', 10000),
('WH-002', 'Secondary Warehouse', '{"city": "Rio de Janeiro", "state": "RJ"}', 5000)
ON CONFLICT (code) DO NOTHING;

-- Permissões finais
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA inventory TO postgres;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA events TO postgres;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA audit TO postgres;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA inventory TO postgres;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA events TO postgres;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA audit TO postgres;

-- Log de inicialização
INSERT INTO audit.audit_logs (action, resource_type, new_values)
VALUES ('DATABASE_INITIALIZED', 'system', '{"status": "completed", "timestamp": "' || NOW() || '"}');

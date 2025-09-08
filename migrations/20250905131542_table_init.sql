-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS wallet (
    wallet_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance INTEGER NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
    version INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wallet_operations (
    id BIGSERIAL PRIMARY KEY,
    wallet_id UUID NOT NULL REFERENCES wallet(wallet_uuid) ON DELETE CASCADE, 
    operation_type VARCHAR(10) NOT NULL CHECK (operation_type IN ('DEPOSIT', 'WITHDRAW')),
    amount INTEGER NOT NULL CHECK (amount > 0),
    reference_id UUID NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPLIED', 'FAILED')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    applied_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_wallet_operations_wallet_id ON wallet_operations(wallet_id);
CREATE INDEX IF NOT EXISTS idx_wallet_operations_reference_id ON wallet_operations(reference_id);
CREATE INDEX IF NOT EXISTS idx_wallet_operations_pending ON wallet_operations(status) WHERE status = 'PENDING';

ALTER TABLE wallet_operations SET (
    autovacuum_enabled = true,
    autovacuum_vacuum_scale_factor = 0.1,
    autovacuum_analyze_scale_factor = 0.05
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_wallet_operations_pending;
DROP INDEX IF EXISTS idx_wallet_operations_reference_id;
DROP INDEX IF EXISTS idx_wallet_operations_wallet_id;

DROP TABLE IF EXISTS wallet_operations;

DROP TABLE IF EXISTS wallet;

DROP EXTENSION IF EXISTS "pgcrypto";
-- +goose StatementEnd

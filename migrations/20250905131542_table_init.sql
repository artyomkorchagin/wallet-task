-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS wallet (
    wallet_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wallet_operations (
    id BIGSERIAL PRIMARY KEY,
    wallet_id UUID NOT NULL REFERENCES wallet(id) ON DELETE CASCADE,
    operation_type VARCHAR(10) NOT NULL CHECK (operation_type IN ('DEPOSIT', 'WITHDRAW')),
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_wallet_operations_wallet_id ON wallet_operations(wallet_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_wallet_operations_wallet_id;

DROP TABLE IF EXISTS wallet_operations;

DROP TABLE IF EXISTS wallet;

DROP EXTENSION IF EXISTS "pgcrypto";
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
INSERT INTO wallet (wallet_uuid, balance, version) VALUES
('a1b2c3e4-5678-9012-3456-789012345678', 1500, 0),
('b2c3d4e5-6789-0123-4567-890123456789', 2750, 1),
('c3d4e5f6-7890-1234-5678-901234567890', 9999, 3),
('d4e5f6a7-8901-2345-6789-012345678901', 0, 0),
('e5f6a7b8-9012-3456-7890-123456789012', 4200, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE * FROM wallet;
-- +goose StatementEnd

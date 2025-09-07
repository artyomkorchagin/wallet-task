package walletpostgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/artyomkorchagin/wallet-task/internal/types"
)

func (r *Repository) UpdateBalance(ctx context.Context, req *types.WalletUpdateRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = r.markOperationFailed(ctx, req.ReferenceID)
		}
		tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, `
        INSERT INTO wallet_operations 
            (wallet_id, operation_type, amount, reference_id, status, created_at)
        VALUES ($1, $2, $3, $4, 'PENDING', NOW())
    `, req.WalletUUID, req.Operation, req.Amount, req.ReferenceID)

	if err != nil {
		return fmt.Errorf("failed to log operation: %w", err)
	}

	var currentBalance, currentVersion int
	err = tx.QueryRowContext(ctx, `
        SELECT balance, version FROM wallet WHERE wallet_uuid = $1 FOR UPDATE
    `, req.WalletUUID).Scan(&currentBalance, &currentVersion)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.ErrWalletNotFound
		}
		return fmt.Errorf("failed to load wallet: %w", err)
	}

	newBalance := currentBalance
	if req.Operation != "DEPOSIT" && req.Operation != "WITHDRAW" {
		return types.ErrInvalidOperation
	}
	if req.Operation == "DEPOSIT" {
		newBalance += req.Amount
	}
	if req.Operation == "WITHDRAW" {
		if currentBalance < req.Amount {
			return types.ErrInsufficientFunds
		}
		newBalance -= req.Amount
	}

	result, err := tx.ExecContext(ctx, `
        UPDATE wallet SET balance = $1, version = version + 1, updated_at = NOW()
        WHERE wallet_uuid = $2 AND version = $3
    `, newBalance, req.WalletUUID, currentVersion)

	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return types.ErrConcurrentUpdate
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE wallet_operations 
        SET status = 'APPLIED', applied_at = NOW() 
        WHERE reference_id = $1
    `, req.ReferenceID)

	if err != nil {
		return fmt.Errorf("failed to mark as applied: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (r *Repository) markOperationFailed(ctx context.Context, referenceID string) error {
	_, err := r.db.ExecContext(ctx, `
        UPDATE wallet_operations 
        SET status = 'FAILED', applied_at = NOW() 
        WHERE reference_id = $1 AND status = 'PENDING'
    `, referenceID)
	return err
}

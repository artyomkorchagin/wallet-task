package walletpostgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/artyomkorchagin/wallet-task/internal/types"
)

func (r *Repository) GetBalance(ctx context.Context, walletUUID string) (int, error) {
	var balance int

	query := "SELECT balance FROM wallet WHERE wallet_uuid = $1"
	err := r.db.QueryRowContext(ctx, query, walletUUID).Scan(&balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, types.ErrWalletNotFound
		}
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

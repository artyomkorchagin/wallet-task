package walletpostgresql

import (
	"context"
	"fmt"
)

func (r *Repository) CheckOperationExists(ctx context.Context, referenceID string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM wallet_operations WHERE reference_id = $1)"

	err := r.db.QueryRowContext(ctx, query, referenceID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check operation existence: %w", err)
	}

	return exists, nil
}

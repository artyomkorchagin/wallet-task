package walletpostgresql

import "context"

func (r *Repository) GetBalance(ctx context.Context, walletUUID string) (int, error) {
	return 0, nil
}

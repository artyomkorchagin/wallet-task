// repository/postgresql/get_balance_test.go
package walletpostgresql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: db}
	ctx := context.Background()
	walletUUID := "a1b2c3e4-5678-9012-3456-789012345678"

	t.Run("success - wallet exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"balance"}).AddRow(500)
		mock.ExpectQuery(`SELECT balance FROM wallet WHERE wallet_uuid = \$1`).
			WithArgs(walletUUID).
			WillReturnRows(rows)

		balance, err := repo.GetBalance(ctx, walletUUID)
		require.NoError(t, err)
		assert.Equal(t, 500, balance)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("wallet not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT balance FROM wallet WHERE wallet_uuid = \$1`).
			WithArgs(walletUUID).
			WillReturnError(sql.ErrNoRows)

		balance, err := repo.GetBalance(ctx, walletUUID)
		require.Error(t, err)
		assert.Equal(t, 0, balance)
		assert.Equal(t, types.ErrWalletNotFound, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT balance FROM wallet WHERE wallet_uuid = \$1`).
			WithArgs(walletUUID).
			WillReturnError(assert.AnError)

		balance, err := repo.GetBalance(ctx, walletUUID)
		require.Error(t, err)
		assert.Equal(t, 0, balance)
		assert.Contains(t, err.Error(), "failed to get balance")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"balance"}).AddRow("not_a_number")
		mock.ExpectQuery(`SELECT balance FROM wallet WHERE wallet_uuid = \$1`).
			WithArgs(walletUUID).
			WillReturnRows(rows)

		balance, err := repo.GetBalance(ctx, walletUUID)
		require.Error(t, err)
		assert.Equal(t, 0, balance)
		assert.Contains(t, err.Error(), "failed to get balance")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

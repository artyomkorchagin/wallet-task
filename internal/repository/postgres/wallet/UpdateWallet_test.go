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

func TestRepository_UpdateBalance(t *testing.T) {
	newValidRequest := func() *types.WalletUpdateRequest {
		return &types.WalletUpdateRequest{
			WalletUUID:  "a1b2c3e4-5678-9012-3456-789012345678",
			Operation:   "DEPOSIT",
			Amount:      100,
			ReferenceID: "ref-123",
		}
	}

	t.Run("success deposit", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{db: db}
		ctx := context.Background()

		req := newValidRequest()

		mock.ExpectBegin()

		mock.ExpectExec(`INSERT INTO wallet_operations`).
			WithArgs(req.WalletUUID, req.Operation, req.Amount, req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		rows := sqlmock.NewRows([]string{"balance", "version"}).AddRow(500, 1)
		mock.ExpectQuery(`SELECT balance, version FROM wallet WHERE wallet_uuid = \$1 FOR UPDATE`).
			WithArgs(req.WalletUUID).
			WillReturnRows(rows)

		mock.ExpectExec(`UPDATE wallet SET balance = \$1, version = version \+ 1, updated_at = NOW\(\) WHERE wallet_uuid = \$2 AND version = \$3`).
			WithArgs(600, req.WalletUUID, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`UPDATE wallet_operations SET status = 'APPLIED', applied_at = NOW\(\) WHERE reference_id = \$1`).
			WithArgs(req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = repo.UpdateBalance(ctx, req)
		require.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success withdraw", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{db: db}
		ctx := context.Background()

		req := newValidRequest()
		req.Operation = "WITHDRAW"
		req.Amount = 50
		req.ReferenceID = "ref-456"

		mock.ExpectBegin()

		mock.ExpectExec(`INSERT INTO wallet_operations`).
			WithArgs(req.WalletUUID, req.Operation, req.Amount, req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		rows := sqlmock.NewRows([]string{"balance", "version"}).AddRow(100, 1)
		mock.ExpectQuery(`SELECT balance, version FROM wallet WHERE wallet_uuid = \$1 FOR UPDATE`).
			WithArgs(req.WalletUUID).
			WillReturnRows(rows)

		mock.ExpectExec(`UPDATE wallet SET balance = \$1, version = version \+ 1, updated_at = NOW\(\) WHERE wallet_uuid = \$2 AND version = \$3`).
			WithArgs(50, req.WalletUUID, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`UPDATE wallet_operations SET status = 'APPLIED', applied_at = NOW\(\) WHERE reference_id = \$1`).
			WithArgs(req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = repo.UpdateBalance(ctx, req)
		require.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("wallet not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{db: db}
		ctx := context.Background()

		req := newValidRequest()

		mock.ExpectBegin()

		mock.ExpectExec(`INSERT INTO wallet_operations`).
			WithArgs(req.WalletUUID, req.Operation, req.Amount, req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`SELECT balance, version FROM wallet WHERE wallet_uuid = \$1 FOR UPDATE`).
			WithArgs(req.WalletUUID).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec(`UPDATE wallet_operations SET status = 'FAILED', applied_at = NOW\(\) WHERE reference_id = \$1 AND status = 'PENDING'`).
			WithArgs(req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectRollback()

		err = repo.UpdateBalance(ctx, req)
		require.Error(t, err)
		assert.Equal(t, types.ErrWalletNotFound, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("insufficient funds", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{db: db}
		ctx := context.Background()

		req := newValidRequest()
		req.Operation = "WITHDRAW"
		req.Amount = 200
		req.ReferenceID = "ref-789"

		mock.ExpectBegin()

		mock.ExpectExec(`INSERT INTO wallet_operations`).
			WithArgs(req.WalletUUID, req.Operation, req.Amount, req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		rows := sqlmock.NewRows([]string{"balance", "version"}).AddRow(100, 1)
		mock.ExpectQuery(`SELECT balance, version FROM wallet WHERE wallet_uuid = \$1 FOR UPDATE`).
			WithArgs(req.WalletUUID).
			WillReturnRows(rows)

		mock.ExpectExec(`UPDATE wallet_operations SET status = 'FAILED', applied_at = NOW\(\) WHERE reference_id = \$1 AND status = 'PENDING'`).
			WithArgs(req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectRollback()

		err = repo.UpdateBalance(ctx, req)
		require.Error(t, err)
		assert.Equal(t, types.ErrInsufficientFunds, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("concurrent update", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{db: db}
		ctx := context.Background()

		req := newValidRequest()

		mock.ExpectBegin()

		mock.ExpectExec(`INSERT INTO wallet_operations`).
			WithArgs(req.WalletUUID, req.Operation, req.Amount, req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		rows := sqlmock.NewRows([]string{"balance", "version"}).AddRow(500, 1)
		mock.ExpectQuery(`SELECT balance, version FROM wallet WHERE wallet_uuid = \$1 FOR UPDATE`).
			WithArgs(req.WalletUUID).
			WillReturnRows(rows)

		mock.ExpectExec(`UPDATE wallet SET balance = \$1, version = version \+ 1, updated_at = NOW\(\) WHERE wallet_uuid = \$2 AND version = \$3`).
			WithArgs(600, req.WalletUUID, 1).
			WillReturnResult(sqlmock.NewResult(1, 0)) // 0 rows affected

		mock.ExpectExec(`UPDATE wallet_operations SET status = 'FAILED', applied_at = NOW\(\) WHERE reference_id = \$1 AND status = 'PENDING'`).
			WithArgs(req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectRollback()

		err = repo.UpdateBalance(ctx, req)
		require.Error(t, err)
		assert.Equal(t, types.ErrConcurrentUpdate, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid operation type", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{db: db}
		ctx := context.Background()

		req := newValidRequest()
		req.Operation = "INVALID"
		req.ReferenceID = "ref-999"

		mock.ExpectBegin()

		mock.ExpectExec(`INSERT INTO wallet_operations`).
			WithArgs(req.WalletUUID, req.Operation, req.Amount, req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		rows := sqlmock.NewRows([]string{"balance", "version"}).AddRow(500, 1)
		mock.ExpectQuery(`SELECT balance, version FROM wallet WHERE wallet_uuid = \$1 FOR UPDATE`).
			WithArgs(req.WalletUUID).
			WillReturnRows(rows)

		mock.ExpectExec(`UPDATE wallet_operations SET status = 'FAILED', applied_at = NOW\(\) WHERE reference_id = \$1 AND status = 'PENDING'`).
			WithArgs(req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectRollback()

		err = repo.UpdateBalance(ctx, req)
		require.Error(t, err)
		assert.Equal(t, types.ErrInvalidOperation, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error on insert", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{db: db}
		ctx := context.Background()

		req := newValidRequest()

		mock.ExpectBegin()

		mock.ExpectExec(`INSERT INTO wallet_operations`).
			WithArgs(req.WalletUUID, req.Operation, req.Amount, req.ReferenceID).
			WillReturnError(assert.AnError)

		mock.ExpectRollback()

		err = repo.UpdateBalance(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to log operation")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error on update wallet", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{db: db}
		ctx := context.Background()

		req := newValidRequest()

		mock.ExpectBegin()

		mock.ExpectExec(`INSERT INTO wallet_operations`).
			WithArgs(req.WalletUUID, req.Operation, req.Amount, req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		rows := sqlmock.NewRows([]string{"balance", "version"}).AddRow(500, 1)
		mock.ExpectQuery(`SELECT balance, version FROM wallet WHERE wallet_uuid = \$1 FOR UPDATE`).
			WithArgs(req.WalletUUID).
			WillReturnRows(rows)

		mock.ExpectExec(`UPDATE wallet SET balance = \$1, version = version \+ 1, updated_at = NOW\(\) WHERE wallet_uuid = \$2 AND version = \$3`).
			WithArgs(600, req.WalletUUID, 1).
			WillReturnError(assert.AnError)

		mock.ExpectExec(`UPDATE wallet_operations SET status = 'FAILED', applied_at = NOW\(\) WHERE reference_id = \$1 AND status = 'PENDING'`).
			WithArgs(req.ReferenceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectRollback()

		err = repo.UpdateBalance(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update wallet")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

package walletpostgresql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_CheckOperationExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: db}

	ctx := context.Background()
	referenceID := "ref-123"

	t.Run("operation exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallet_operations WHERE reference_id = \$1\)`).
			WithArgs(referenceID).
			WillReturnRows(rows)

		exists, err := repo.CheckOperationExists(ctx, referenceID)
		require.NoError(t, err)
		assert.True(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("operation does not exist", func(t *testing.T) {
		// Мок: возвращаем false
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallet_operations WHERE reference_id = \$1\)`).
			WithArgs(referenceID).
			WillReturnRows(rows)

		exists, err := repo.CheckOperationExists(ctx, referenceID)
		require.NoError(t, err)
		assert.False(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallet_operations WHERE reference_id = \$1\)`).
			WithArgs(referenceID).
			WillReturnError(assert.AnError)

		exists, err := repo.CheckOperationExists(ctx, referenceID)
		require.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), "failed to check operation existence")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

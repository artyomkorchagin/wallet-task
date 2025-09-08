package walletservice_test

import (
	"context"
	"testing"

	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateBalance_Validation(t *testing.T) {
	service, _ := setupService(t)

	ctx := context.Background()
	validUUID := uuid.New().String()

	t.Run("empty walletUUID", func(t *testing.T) {
		req := &types.WalletUpdateRequest{
			WalletUUID:  "",
			Operation:   types.OperationTypeDeposit,
			Amount:      100,
			ReferenceID: validUUID,
		}

		err := service.UpdateBalance(ctx, req)
		require.Error(t, err)
		assert.Equal(t, 400, err.(types.HTTPError).Code)
		assert.Contains(t, err.Error(), "walletUUID is empty")
	})

}

func TestUpdateBalance_Idempotency(t *testing.T) {
	service, repo := setupService(t)
	ctx := context.Background()
	refID := uuid.New().String()
	wur := &types.WalletUpdateRequest{
		WalletUUID:  uuid.New().String(),
		Operation:   types.OperationTypeDeposit,
		Amount:      100,
		ReferenceID: refID,
	}

	t.Run("operation already exists", func(t *testing.T) {
		repo.On("CheckOperationExists", ctx, refID).Return(true, nil)

		err := service.UpdateBalance(ctx, wur)
		require.Error(t, err)
		assert.Equal(t, 409, err.(types.HTTPError).Code)
		assert.Contains(t, err.Error(), "already processed")

		repo.AssertExpectations(t)
	})
}

func TestUpdateBalance_RepositoryErrors(t *testing.T) {
	service, repo := setupService(t)

	ctx := context.Background()
	req := &types.WalletUpdateRequest{
		WalletUUID:  uuid.New().String(),
		Operation:   types.OperationTypeWithdraw,
		Amount:      100,
		ReferenceID: uuid.New().String(),
	}

	t.Run("wallet not found", func(t *testing.T) {
		repo.On("CheckOperationExists", ctx, mock.Anything).Return(false, nil)
		repo.On("UpdateBalance", ctx, req).Return(types.ErrWalletNotFound)

		err := service.UpdateBalance(ctx, req)
		require.Error(t, err)
		assert.Equal(t, 404, err.(types.HTTPError).Code)
		assert.Contains(t, err.Error(), "wallet not found")

		repo.AssertExpectations(t)
	})

}

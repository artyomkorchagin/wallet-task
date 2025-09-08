package walletservice_test

import (
	"context"

	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CheckOperationExists(ctx context.Context, referenceID string) (bool, error) {
	args := m.Called(ctx, referenceID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) UpdateBalance(ctx context.Context, req *types.WalletUpdateRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockRepository) GetBalance(ctx context.Context, walletUUID string) (int, error) {
	args := m.Called(ctx, walletUUID)
	return args.Int(0), args.Error(1)
}

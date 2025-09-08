package router

import (
	"context"

	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/stretchr/testify/mock"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) GetBalance(ctx context.Context, walletUUID string) (int, error) {
	args := m.Called(ctx, walletUUID)
	return args.Int(0), args.Error(1)
}

func (m *MockWalletService) UpdateBalance(ctx context.Context, req *types.WalletUpdateRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

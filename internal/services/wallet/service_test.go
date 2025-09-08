package walletservice_test

import (
	"context"
	"testing"

	walletservice "github.com/artyomkorchagin/wallet-task/internal/services/wallet"
	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
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

func setupService(t *testing.T) (*walletservice.Service, *MockRepository) {
	redisMockClient, _ := redismock.NewClientMock()
	logger := zaptest.NewLogger(t)
	repo := new(MockRepository)
	service := walletservice.NewService(repo, redisMockClient, logger)
	return service, repo
}

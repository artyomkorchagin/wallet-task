package walletservice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	walletservice "github.com/artyomkorchagin/wallet-task/internal/services/wallet"
	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func setupService(t *testing.T) (*walletservice.Service, *MockRepository) {
	redisMockClient, _ := redismock.NewClientMock()
	logger := zaptest.NewLogger(t)
	repo := new(MockRepository)
	service := walletservice.NewService(repo, redisMockClient, logger)
	return service, repo
}

func TestGetBalance_Validation(t *testing.T) {

	service, _ := setupService(t)

	ctx := context.Background()

	t.Run("empty walletUUID", func(t *testing.T) {
		balance, err := service.GetBalance(ctx, "")
		require.Error(t, err)
		assert.Equal(t, 400, err.(types.HTTPError).Code)
		assert.Contains(t, err.Error(), "walletUUID is empty")
		assert.Equal(t, 0, balance)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		balance, err := service.GetBalance(ctx, "invalid-uuid")
		require.Error(t, err)
		assert.Equal(t, 400, err.(types.HTTPError).Code)
		assert.Contains(t, err.Error(), "not valid")
		assert.Equal(t, 0, balance)
	})
}

func TestGetBalance_CacheHit(t *testing.T) {
	redisMockClient, redisMock := redismock.NewClientMock()
	logger := zaptest.NewLogger(t)
	repo := new(MockRepository)
	service := walletservice.NewService(repo, redisMockClient, logger)

	ctx := context.Background()
	walletUUID := uuid.New().String()
	cacheKey := "wallet:balance:" + walletUUID

	redisMock.ExpectGet(cacheKey).SetVal("500")

	balance, err := service.GetBalance(ctx, walletUUID)
	require.NoError(t, err)
	assert.Equal(t, 500, balance)

	repo.AssertNotCalled(t, "GetBalance", mock.Anything, walletUUID)
	redisMock.ExpectationsWereMet()
}

func TestGetBalance_CacheMiss(t *testing.T) {
	redisMockClient, redisMock := redismock.NewClientMock()
	logger := zaptest.NewLogger(t)
	repo := new(MockRepository)
	service := walletservice.NewService(repo, redisMockClient, logger)

	ctx := context.Background()
	walletUUID := uuid.New().String()
	cacheKey := "wallet:balance:" + walletUUID

	redisMock.ExpectGet(cacheKey).RedisNil()

	repo.On("GetBalance", ctx, walletUUID).Return(300, nil)

	redisMock.ExpectSet(cacheKey, "300", 15*time.Second).SetVal("OK")

	balance, err := service.GetBalance(ctx, walletUUID)
	require.NoError(t, err)
	assert.Equal(t, 300, balance)

	repo.AssertExpectations(t)
	redisMock.ExpectationsWereMet()
}

func TestGetBalance_RepositoryNotFound(t *testing.T) {
	redisMockClient, redisMock := redismock.NewClientMock()
	logger := zaptest.NewLogger(t)
	repo := new(MockRepository)
	service := walletservice.NewService(repo, redisMockClient, logger)

	ctx := context.Background()
	walletUUID := uuid.New().String()
	cacheKey := "wallet:balance:" + walletUUID

	redisMock.ExpectGet(cacheKey).RedisNil()
	repo.On("GetBalance", ctx, walletUUID).Return(0, types.ErrWalletNotFound)

	balance, err := service.GetBalance(ctx, walletUUID)
	require.Error(t, err)
	assert.Equal(t, 404, err.(types.HTTPError).Code)
	assert.Contains(t, err.Error(), "wallet not found")
	assert.Equal(t, 0, balance)

	repo.AssertExpectations(t)
	redisMock.ExpectationsWereMet()
}

func TestGetBalance_RepositoryError(t *testing.T) {
	redisMockClient, redisMock := redismock.NewClientMock()
	logger := zaptest.NewLogger(t)
	repo := new(MockRepository)
	service := walletservice.NewService(repo, redisMockClient, logger)

	ctx := context.Background()
	walletUUID := uuid.New().String()
	cacheKey := "wallet:balance:" + walletUUID

	redisMock.ExpectGet(cacheKey).RedisNil()

	dbErr := errors.New("pq: server unavailable")
	repo.On("GetBalance", ctx, walletUUID).Return(0, dbErr)

	balance, err := service.GetBalance(ctx, walletUUID)

	require.Error(t, err)
	assert.Equal(t, 0, balance)

	httpErr, ok := err.(types.HTTPError)
	require.True(t, ok, "error should be types.HTTPError")
	assert.Equal(t, 500, httpErr.Code)

	assert.Contains(t, err.Error(), "pq: server unavailable")
}

func TestGetBalance_CacheSetFails(t *testing.T) {
	redisMockClient, redisMock := redismock.NewClientMock()
	logger := zaptest.NewLogger(t)
	repo := new(MockRepository)
	service := walletservice.NewService(repo, redisMockClient, logger)

	ctx := context.Background()
	walletUUID := uuid.New().String()
	cacheKey := "wallet:balance:" + walletUUID

	redisMock.ExpectGet(cacheKey).RedisNil()
	repo.On("GetBalance", ctx, walletUUID).Return(200, nil)

	redisMock.ExpectSet(cacheKey, "200", 15*time.Second).SetErr(errors.New("redis down"))

	balance, err := service.GetBalance(ctx, walletUUID)
	require.NoError(t, err)
	assert.Equal(t, 200, balance)

	repo.AssertExpectations(t)
	redisMock.ExpectationsWereMet()
}

package walletservice

import (
	"context"

	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ServiceInterface interface {
	GetBalance(ctx context.Context, walletUUID string) (int, error)
	UpdateBalance(ctx context.Context, req *types.WalletUpdateRequest) error
}

type Service struct {
	repo   ReadWriter
	redis  *redis.Client
	logger *zap.Logger
}

func NewService(repo ReadWriter, redis *redis.Client, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		redis:  redis,
		logger: logger,
	}
}

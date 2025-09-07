package walletservice

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

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

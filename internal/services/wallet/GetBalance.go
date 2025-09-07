package walletservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) GetBalance(ctx context.Context, walletUUID string) (int, error) {
	start := time.Now()
	logger := s.logger.With(zap.String("wallet_uuid", walletUUID))
	defer func() {
		duration := time.Since(start)
		if duration > 100*time.Millisecond {
			logger.Warn("GetBalance slow",
				zap.Duration("duration", duration))
		} else {
			logger.Debug("GetBalance finished",
				zap.Duration("duration", duration))
		}
	}()
	logger.Info("GetBalance called")

	if walletUUID == "" {
		logger.Warn("GetBalance: walletUUID is empty")
		return 0, types.ErrBadRequest(fmt.Errorf("walletUUID is empty"))
	}

	_, err := uuid.Parse(walletUUID)
	if err != nil {
		logger.Warn("GetBalance: invalid UUID",
			zap.Error(err))
		return 0, types.ErrBadRequest(fmt.Errorf("walletUUID is not valid: %w", err))
	}

	cacheKey := "wallet:balance:" + walletUUID

	if balance, err := s.redis.Get(ctx, cacheKey).Int(); err == nil {
		logger.Debug("GetBalance: cache hit",
			zap.Int("balance", balance))
		return balance, nil
	}

	logger.Debug("GetBalance: cache miss",
		zap.String("cache_key", cacheKey))

	balance, err := s.repo.GetBalance(ctx, walletUUID)
	if err != nil {
		logger.Error("GetBalance: failed to load from repository",
			zap.Error(err))

		if errors.Is(err, types.ErrWalletNotFound) {
			return 0, types.ErrNotFound(err)
		}
		return 0, types.ErrInternalServerError(err)
	}

	if setErr := s.redis.Set(ctx, cacheKey, balance, 15*time.Second).Err(); setErr != nil {
		logger.Warn("GetBalance: failed to set cache",
			zap.String("cache_key", cacheKey),
			zap.Error(setErr))
		// Не возвращаем ошибку, кеш опционален
	} else {
		logger.Debug("GetBalance: balance cached",
			zap.Int("balance", balance),
			zap.Duration("ttl", 15*time.Second))
	}

	logger.Info("GetBalance: success",
		zap.Int("balance", balance))

	return balance, nil
}

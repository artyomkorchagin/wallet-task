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

func (s *Service) UpdateBalance(ctx context.Context, wur *types.WalletUpdateRequest) error {
	start := time.Now()
	logger := s.logger.With(zap.String("wallet_uuid", wur.WalletUUID))
	defer func() {
		duration := time.Since(start)
		if duration > 100*time.Millisecond {
			logger.Warn("UpdateBalance slow",
				zap.Duration("duration", duration))
		} else {
			logger.Debug("UpdateBalance finished",
				zap.Duration("duration", duration))
		}
	}()
	logger.Info("UpdateBalance called",
		zap.String("operation", wur.Operation),
		zap.Int("amount", wur.Amount),
		zap.String("reference_id", wur.ReferenceID),
	)

	if wur.WalletUUID == "" {
		logger.Warn("UpdateBalance: walletUUID is empty")
		return types.ErrBadRequest(fmt.Errorf("walletUUID is empty"))
	}

	if wur.Amount <= 0 {
		logger.Warn("UpdateBalance: invalid amount",
			zap.Int("amount", wur.Amount))
		return types.ErrBadRequest(fmt.Errorf("amount must be positive"))
	}

	if wur.Operation != types.OperationTypeDeposit && wur.Operation != types.OperationTypeWithdraw {
		logger.Warn("UpdateBalance: invalid operation type",
			zap.String("operation", wur.Operation))
		return types.ErrBadRequest(fmt.Errorf(
			"operation must be either %s or %s",
			types.OperationTypeDeposit,
			types.OperationTypeWithdraw,
		))
	}

	if _, err := uuid.Parse(wur.ReferenceID); err != nil {
		logger.Warn("UpdateBalance: invalid reference_id",
			zap.String("reference_id", wur.ReferenceID),
			zap.Error(err))
		return types.ErrBadRequest(fmt.Errorf("ReferenceID is not valid UUID: %w", err))
	}

	if _, err := uuid.Parse(wur.WalletUUID); err != nil {
		logger.Warn("UpdateBalance: invalid wallet_uuid", zap.Error(err))
		return types.ErrBadRequest(fmt.Errorf("WalletUUID is not valid UUID: %w", err))
	}

	logger.Debug("UpdateBalance: checking idempotency",
		zap.String("reference_id", wur.ReferenceID))

	exists, err := s.repo.CheckOperationExists(ctx, wur.ReferenceID)
	if err != nil {
		logger.Error("UpdateBalance: failed to check idempotency",
			zap.String("reference_id", wur.ReferenceID),
			zap.Error(err))
		return types.ErrInternalServerError(fmt.Errorf("failed to check idempotency: %w", err))
	}

	if exists {
		logger.Warn("UpdateBalance: idempotency conflict")
		return types.ErrConflict(fmt.Errorf("operation with reference_id %s already processed", wur.ReferenceID))
	}

	logger.Debug("UpdateBalance: idempotency check passed",
		zap.String("reference_id", wur.ReferenceID))

	if err = s.repo.UpdateBalance(ctx, wur); err != nil {
		switch {
		case errors.Is(err, types.ErrWalletNotFound):
			logger.Info("UpdateBalance: wallet not found")
			return types.ErrNotFound(err)
		case errors.Is(err, types.ErrInsufficientFunds):
			logger.Info("UpdateBalance: insufficient funds",
				zap.Int("amount", wur.Amount))
			return types.ErrBadRequest(err)
		case errors.Is(err, types.ErrConcurrentUpdate):
			logger.Warn("UpdateBalance: concurrent update",
				zap.String("reference_id", wur.ReferenceID))
			return types.ErrConflict(err)
		case errors.Is(err, types.ErrOperationExists):
			logger.Warn("UpdateBalance: operation exists (race)",
				zap.String("reference_id", wur.ReferenceID))
			return types.ErrConflict(err)
		default:
			logger.Error("UpdateBalance: unexpected error",
				zap.String("operation", wur.Operation),
				zap.Int("amount", wur.Amount),
				zap.Error(err))
			return types.ErrInternalServerError(err)
		}
	}

	logger.Info("UpdateBalance: success",
		zap.String("operation", wur.Operation),
		zap.Int("amount", wur.Amount),
		zap.String("reference_id", wur.ReferenceID))

	return nil
}

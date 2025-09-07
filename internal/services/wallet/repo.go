package walletservice

import (
	"context"

	"github.com/artyomkorchagin/wallet-task/internal/types"
)

type Reader interface {
	GetBalance(ctx context.Context, walletUUID string) (int, error)
	CheckOperationExists(ctx context.Context, referenceID string) (bool, error)
}

type Writer interface {
	UpdateBalance(ctx context.Context, wallet *types.WalletUpdateRequest) error
}

type ReadWriter interface {
	Reader
	Writer
}

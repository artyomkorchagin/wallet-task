package types

import (
	"github.com/google/uuid"
)

type WalletUpdateRequest struct {
	WalletUUID    uuid.UUID `json:"wallet_uuid"`
	OperationType string    `json:"operation_type"`
	Amount        int       `json:"amount"`
}

func NewWalletUpdateRequest(walletUUID uuid.UUID, operationType string, amount int) *WalletUpdateRequest {
	return &WalletUpdateRequest{
		walletUUID,
		operationType,
		amount,
	}
}

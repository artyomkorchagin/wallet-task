package types

type WalletUpdateRequest struct {
	WalletUUID  string `json:"valletId" binding:"required,uuid4"`
	Operation   string `json:"operationType" binding:"required,oneof=DEPOSIT WITHDRAW"`
	Amount      int    `json:"amount" binding:"required,gt=0"`
	ReferenceID string
}

func NewWalletUpdateRequest(walletUUID string, operationType string, amount int, referenceid string) *WalletUpdateRequest {
	return &WalletUpdateRequest{
		walletUUID,
		operationType,
		amount,
		referenceid,
	}
}

var (
	OperationTypeDeposit  = "DEPOSIT"
	OperationTypeWithdraw = "WITHDRAW"
)

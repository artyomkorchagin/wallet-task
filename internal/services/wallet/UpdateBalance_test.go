package walletservice

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/artyomkorchagin/wallet-task/internal/services/wallet/mocks"
	"github.com/artyomkorchagin/wallet-task/internal/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestService_UpdateBalance(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.WalletUpdateRequest
		mockSetup   func(*mocks.MockReadWriter)
		expectedErr error
	}{
		{
			name: "success deposit",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup: func(m *mocks.MockReadWriter) {
				m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, nil)
				m.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "success withdraw",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeWithdraw,
				Amount:      50,
				ReferenceID: uuid.New().String(),
			},
			mockSetup: func(m *mocks.MockReadWriter) {
				m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, nil)
				m.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "empty walletUUID",
			req: &types.WalletUpdateRequest{
				WalletUUID:  "",
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup:   func(m *mocks.MockReadWriter) {},
			expectedErr: types.ErrBadRequest(fmt.Errorf("walletUUID is empty")),
		},
		{
			name: "invalid amount (zero)",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeDeposit,
				Amount:      0,
				ReferenceID: uuid.New().String(),
			},
			mockSetup:   func(m *mocks.MockReadWriter) {},
			expectedErr: types.ErrBadRequest(fmt.Errorf("amount must be positive")),
		},
		{
			name: "invalid operation",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   "INVALID",
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup:   func(m *mocks.MockReadWriter) {},
			expectedErr: types.ErrBadRequest(fmt.Errorf("operation must be either DEPOSIT or WITHDRAW")),
		},
		{
			name: "invalid reference_id UUID",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: "not-a-uuid",
			},
			mockSetup:   func(m *mocks.MockReadWriter) {},
			expectedErr: types.ErrBadRequest(fmt.Errorf("ReferenceID is not valid UUID")),
		},
		{
			name: "invalid wallet_uuid UUID",
			req: &types.WalletUpdateRequest{
				WalletUUID:  "not-a-uuid",
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup:   func(m *mocks.MockReadWriter) {},
			expectedErr: types.ErrBadRequest(fmt.Errorf("WalletUUID is not valid UUID")),
		},
		{
			name: "idempotency conflict",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup: func(m *mocks.MockReadWriter) {
				m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(true, nil)
			},
			expectedErr: types.ErrConflict(errors.New("operation with reference_id")),
		},
		{
			name: "CheckOperationExists internal error",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup: func(m *mocks.MockReadWriter) {
				m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, errors.New("db error"))
			},
			expectedErr: types.ErrInternalServerError(fmt.Errorf("failed to check idempotency")),
		},
		{
			name: "wallet not found",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup: func(m *mocks.MockReadWriter) {
				m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, nil)
				m.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).Return(types.ErrWalletNotFound)
			},
			expectedErr: types.ErrNotFound(types.ErrWalletNotFound),
		},
		{
			name: "insufficient funds",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeWithdraw,
				Amount:      1000,
				ReferenceID: uuid.New().String(),
			},
			mockSetup: func(m *mocks.MockReadWriter) {
				m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, nil)
				m.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).Return(types.ErrInsufficientFunds)
			},
			expectedErr: types.ErrBadRequest(types.ErrInsufficientFunds),
		},
		{
			name: "concurrent update (conflict)",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup: func(m *mocks.MockReadWriter) {
				m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, nil)
				m.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).Return(types.ErrConcurrentUpdate)
			},
			expectedErr: types.ErrConflict(types.ErrConcurrentUpdate),
		},
		{
			name: "operation exists race",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup: func(m *mocks.MockReadWriter) {
				m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, nil)
				m.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).Return(types.ErrOperationExists)
			},
			expectedErr: types.ErrConflict(types.ErrOperationExists),
		},
		{
			name: "unexpected repo error",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup: func(m *mocks.MockReadWriter) {
				m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, nil)
				m.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).Return(errors.New("unknown db error"))
			},
			expectedErr: types.ErrInternalServerError(errors.New("unknown db error")),
		},
		{
			name: "deadlock retry success on 2nd attempt",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup: func(m *mocks.MockReadWriter) {
				gomock.InOrder(
					m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, nil),
					m.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).Return(errors.New("deadlock detected")),
					m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, nil),
					m.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			expectedErr: nil,
		},
		{
			name: "deadlock retry fails after 3 attempts",
			req: &types.WalletUpdateRequest{
				WalletUUID:  uuid.New().String(),
				Operation:   types.OperationTypeDeposit,
				Amount:      100,
				ReferenceID: uuid.New().String(),
			},
			mockSetup: func(m *mocks.MockReadWriter) {
				for i := 0; i < 3; i++ {
					m.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, nil).Times(1)
					m.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).Return(errors.New("deadlock detected")).Times(1)
				}
			},
			expectedErr: types.ErrConflict(fmt.Errorf("too many retries due to deadlock")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockReadWriter(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			logger, _ := zap.NewDevelopment()

			svc := &Service{
				repo:   mockRepo,
				redis:  nil,
				logger: logger,
			}

			err := svc.UpdateBalance(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_UpdateBalance_SlowOperationLogged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockReadWriter(ctrl)
	mockRepo.EXPECT().CheckOperationExists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockRepo.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, req *types.WalletUpdateRequest) error {
		time.Sleep(150 * time.Millisecond)
		return nil
	})

	logger, _ := zap.NewDevelopment()

	svc := &Service{
		repo:   mockRepo,
		redis:  nil,
		logger: logger,
	}

	req := &types.WalletUpdateRequest{
		WalletUUID:  uuid.New().String(),
		Operation:   types.OperationTypeDeposit,
		Amount:      100,
		ReferenceID: uuid.New().String(),
	}

	err := svc.UpdateBalance(context.Background(), req)
	assert.NoError(t, err)
}

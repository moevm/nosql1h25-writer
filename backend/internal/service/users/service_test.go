package users_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"

	log "github.com/sirupsen/logrus"

	users_repo "github.com/moevm/nosql1h25-writer/backend/internal/repo/users"
	users_repo_mocks "github.com/moevm/nosql1h25-writer/backend/internal/repo/users/mocks"
	users_service "github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

func TestService_UpdateBalance_Deposit(t *testing.T) {
	log.SetOutput(io.Discard)

	var (
		ctx         = context.TODO()
		userID      = primitive.NewObjectID()
		amount      = 100
		operation   = users_service.OperationTypeDeposit
	)

	type MockBehavior func(r *users_repo_mocks.MockRepo)

	tests := []struct {
		name         string
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "successful deposit",
			mockBehavior: func(r *users_repo_mocks.MockRepo) {
				r.EXPECT().Deposit(ctx, userID, amount).Return(nil)
			},
		},
		{
			name: "deposit failure",
			mockBehavior: func(r *users_repo_mocks.MockRepo) {
				r.EXPECT().Deposit(ctx, userID, amount).Return(users_repo.ErrCannotDeposit)
			},
			wantErr: users_service.ErrCannotDeposit,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := users_repo_mocks.NewMockRepo(ctrl)
			tc.mockBehavior(mockRepo)

			svc := users_service.New(mockRepo)
			err := svc.UpdateBalance(ctx, userID, operation, amount)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestService_UpdateBalance_Withdraw(t *testing.T) {
	log.SetOutput(io.Discard)

	var (
		ctx         = context.TODO()
		userID      = primitive.NewObjectID()
		validAmount = 50
		failAmount  = 100
	)

	type MockBehavior func(r *users_repo_mocks.MockRepo)

	tests := []struct {
		name         string
		amount       int
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name:   "successful withdraw",
			amount: validAmount,
			mockBehavior: func(r *users_repo_mocks.MockRepo) {
				r.EXPECT().Withdraw(ctx, userID, validAmount).Return(nil)
			},
		},
		{
			name:   "insufficient funds",
			amount: failAmount,
			mockBehavior: func(r *users_repo_mocks.MockRepo) {
				r.EXPECT().Withdraw(ctx, userID, failAmount).Return(users_repo.ErrInsufficientFunds)
			},
			wantErr: users_service.ErrInsufficientFunds,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := users_repo_mocks.NewMockRepo(ctrl)
			tc.mockBehavior(mockRepo)

			svc := users_service.New(mockRepo)
			err := svc.UpdateBalance(ctx, userID, users_service.OperationTypeWithdraw, tc.amount)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

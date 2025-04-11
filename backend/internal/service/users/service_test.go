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
		ctx       = context.TODO()
		userID    = primitive.NewObjectID()
		amount    = 100
		operation = users_service.OperationTypeDeposit
	)

	type MockBehavior func(r *users_repo_mocks.MockRepo)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         int
		wantErr      error
	}{
		{
			name: "successful deposit",
			mockBehavior: func(r *users_repo_mocks.MockRepo) {
				r.EXPECT().Deposit(ctx, userID, amount).Return(amount, nil)
			},
			want: amount,
		},
		{
			name: "deposit failure",
			mockBehavior: func(r *users_repo_mocks.MockRepo) {
				r.EXPECT().Deposit(ctx, userID, amount).Return(0, assert.AnError)
			},
			wantErr: users_service.ErrUpdateBalance,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockRepo := users_repo_mocks.NewMockRepo(ctrl)

			tc.mockBehavior(mockRepo)

			svc := users_service.New(mockRepo)

			got, err := svc.UpdateBalance(ctx, userID, operation, amount)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestService_UpdateBalance_Withdraw(t *testing.T) {
	log.SetOutput(io.Discard)

	var (
		ctx       = context.TODO()
		userID    = primitive.NewObjectID()
		amount    = 900
		operation = users_service.OperationTypeWithdraw
	)

	type MockBehavior func(r *users_repo_mocks.MockRepo)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         int
		wantErr      error
	}{
		{
			name: "successful withdraw",
			mockBehavior: func(r *users_repo_mocks.MockRepo) {
				r.EXPECT().Withdraw(ctx, userID, amount).Return(amount, nil)
			},
			want: amount,
		},
		{
			name: "insufficient funds",
			mockBehavior: func(r *users_repo_mocks.MockRepo) {
				r.EXPECT().Withdraw(ctx, userID, amount).Return(0, users_repo.ErrInsufficientFunds)
			},
			wantErr: users_service.ErrInsufficientFunds,
		},
		{
			name: "unexpected error",
			mockBehavior: func(r *users_repo_mocks.MockRepo) {
				r.EXPECT().Withdraw(ctx, userID, amount).Return(0, assert.AnError)
			},
			wantErr: users_service.ErrUpdateBalance,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockRepo := users_repo_mocks.NewMockRepo(ctrl)

			tc.mockBehavior(mockRepo)

			svc := users_service.New(mockRepo)

			got, err := svc.UpdateBalance(ctx, userID, operation, amount)

			assert.ErrorIs(t, err, tc.wantErr)
			if err == nil {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

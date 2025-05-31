package users_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"

	log "github.com/sirupsen/logrus"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	orders_repo_mocks "github.com/moevm/nosql1h25-writer/backend/internal/repo/orders/mocks"
	users_repo "github.com/moevm/nosql1h25-writer/backend/internal/repo/users"
	users_repo_mocks "github.com/moevm/nosql1h25-writer/backend/internal/repo/users/mocks"
	users_service "github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

func TestService_GetByIDExt(t *testing.T) {
	log.SetOutput(io.Discard)
	var (
		ctx    = context.Background()
		userID = primitive.NewObjectID()
		ctrl   = gomock.NewController(t)
	)
	type MockBehavior func(r *users_repo_mocks.MockRepo)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         entity.UserExt
		expectedErr  error
	}{
		{
			name: "Success - user found",
			mockBehavior: func(m *users_repo_mocks.MockRepo) {
				m.EXPECT().
					GetByIDExt(ctx, userID).
					Return(entity.UserExt{
						User: entity.User{
							ID:    userID,
							Email: "test@example.com",
						},
					}, nil)
			},
			want: entity.UserExt{
				User: entity.User{
					ID:    userID,
					Email: "test@example.com",
				},
			},
		}, {
			name: "Error - user not found",
			mockBehavior: func(m *users_repo_mocks.MockRepo) {
				m.EXPECT().
					GetByIDExt(ctx, userID).
					Return(entity.UserExt{}, users_repo.ErrUserNotFound)
			},
			expectedErr: users_service.ErrUserNotFound,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := users_repo_mocks.NewMockRepo(ctrl)
			tc.mockBehavior(mockRepo)

			service := users_service.New(mockRepo, nil)
			got, err := service.GetByIDExt(ctx, userID)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

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

			svc := users_service.New(mockRepo, nil)

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

			svc := users_service.New(mockRepo, nil)

			got, err := svc.UpdateBalance(ctx, userID, operation, amount)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestService_FindOrdersByUserID(t *testing.T) {
	var (
		ctx      = context.TODO()
		targetID = primitive.NewObjectID()
		orders   = []entity.OrderExt{
			{
				Order: entity.Order{
					ID:       primitive.NewObjectID(),
					ClientID: targetID,
					Title:    "Test Order",
				},
			},
		}
	)

	type MockBehavior func(repo *orders_repo_mocks.MockRepo, usersRepo *users_repo_mocks.MockRepo)

	tests := []struct {
		name         string
		requesterID  primitive.ObjectID
		targetID     primitive.ObjectID
		isAdmin      bool
		mockBehavior MockBehavior
		want         []entity.OrderExt
		wantErr      error
	}{
		{
			name:        "admin can view any user orders",
			requesterID: primitive.NewObjectID(),
			targetID:    targetID,
			isAdmin:     true,
			mockBehavior: func(repo *orders_repo_mocks.MockRepo, usersRepo *users_repo_mocks.MockRepo) {
				usersRepo.EXPECT().GetByIDExt(ctx, targetID).Return(entity.UserExt{}, nil)
				repo.EXPECT().FindByUserIDExt(ctx, targetID, true).Return(orders, nil)
			},
			want: orders,
		},
		{
			name:        "user can view own orders",
			requesterID: targetID,
			targetID:    targetID,
			isAdmin:     false,
			mockBehavior: func(repo *orders_repo_mocks.MockRepo, usersRepo *users_repo_mocks.MockRepo) {
				usersRepo.EXPECT().GetByIDExt(ctx, targetID).Return(entity.UserExt{}, nil)
				repo.EXPECT().FindByUserIDExt(ctx, targetID, false).Return(orders, nil)
			},
			want: orders,
		},
		{
			name:        "user cannot view other user orders",
			requesterID: primitive.NewObjectID(),
			targetID:    targetID,
			isAdmin:     false,
			mockBehavior: func(repo *orders_repo_mocks.MockRepo, usersRepo *users_repo_mocks.MockRepo) {
			},
			wantErr: users_service.ErrForbidden,
		},
		{
			name:        "target user not found",
			requesterID: targetID,
			targetID:    targetID,
			isAdmin:     false,
			mockBehavior: func(repo *orders_repo_mocks.MockRepo, usersRepo *users_repo_mocks.MockRepo) {
				usersRepo.EXPECT().GetByIDExt(ctx, targetID).Return(entity.UserExt{}, users_repo.ErrUserNotFound)
			},
			wantErr: users_service.ErrUserNotFound,
		},
		{
			name:        "error getting orders",
			requesterID: targetID,
			targetID:    targetID,
			isAdmin:     false,
			mockBehavior: func(repo *orders_repo_mocks.MockRepo, usersRepo *users_repo_mocks.MockRepo) {
				usersRepo.EXPECT().GetByIDExt(ctx, targetID).Return(entity.UserExt{}, nil)
				repo.EXPECT().FindByUserIDExt(ctx, targetID, false).Return(nil, assert.AnError)
			},
			wantErr: users_service.ErrCannotFindOrders,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockOrdersRepo := orders_repo_mocks.NewMockRepo(ctrl)
			mockUsersRepo := users_repo_mocks.NewMockRepo(ctrl)

			tc.mockBehavior(mockOrdersRepo, mockUsersRepo)

			s := users_service.New(mockUsersRepo, mockOrdersRepo)

			got, err := s.FindOrdersByUserID(ctx, tc.requesterID, tc.targetID, tc.isAdmin)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

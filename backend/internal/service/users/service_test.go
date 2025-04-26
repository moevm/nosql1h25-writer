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
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)
func TestService_GetUserByID(t *testing.T) {
    log.SetOutput(io.Discard)
    ctx := context.Background()
    userID := primitive.NewObjectID()

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    tests := []struct {
        name        string
        mockSetup   func(*users_repo_mocks.MockRepo)
        wantUser    entity.UserExt
        wantErr     bool
        expectedErr error
    }{
        {
            name: "Success - user found",
            mockSetup: func(m *users_repo_mocks.MockRepo) {
                m.EXPECT().
                    GetByIDExt(ctx, userID).
                    Return(entity.UserExt{
                        User: entity.User{
                            ID: userID,
                            Email: "test@example.com",
                        },
                    }, nil)
            },
            wantUser: entity.UserExt{
                User: entity.User{
                    ID: userID,
                    Email: "test@example.com",
                },
            },
        },
        {
            name: "Error - user not found",
            mockSetup: func(m *users_repo_mocks.MockRepo) {
                m.EXPECT().
                    GetByIDExt(ctx, userID).
                    Return(entity.UserExt{}, users_repo.ErrUserNotFound)
            },
            wantErr:     true,
            expectedErr: users_service.ErrUserNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := users_repo_mocks.NewMockRepo(ctrl)
            tt.mockSetup(mockRepo)

            service := users_service.New(mockRepo)
            got, err := service.GetUserByID(ctx, userID)

            if tt.wantErr {
                assert.Error(t, err)
                assert.ErrorIs(t, err, tt.expectedErr)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantUser, got)
            }
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
			assert.Equal(t, tc.want, got)
		})
	}
}

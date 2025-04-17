package orders_test

import (
	"context"
	"io"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"

	orders_repo "github.com/moevm/nosql1h25-writer/backend/internal/repo/orders"
	orders_repo_mocks "github.com/moevm/nosql1h25-writer/backend/internal/repo/orders/mocks"
	orders_service "github.com/moevm/nosql1h25-writer/backend/internal/service/orders"
)

func TestService_Create(t *testing.T) {
	log.SetOutput(io.Discard)
	var (
		ctx       = context.TODO()
		serviceIn = orders_service.CreateIn{
			ClientID:       primitive.NewObjectID(),
			Title:          "bombordiro_crocodilo",
			Description:    "tralalelo tralala",
			Cost:           100,
			CompletionTime: 3600000000000,
		}
		orderID = primitive.NewObjectID()
	)

	type MockBehavior func(o *orders_repo_mocks.MockRepo)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         primitive.ObjectID
		wantErr      error
	}{
		{
			name: "successful",
			mockBehavior: func(o *orders_repo_mocks.MockRepo) {
				o.EXPECT().Create(ctx, orders_repo.CreateIn(serviceIn)).Return(orderID, nil)
			},
			want: orderID,
		},
		{
			name: "cannot create order",
			mockBehavior: func(o *orders_repo_mocks.MockRepo) {
				o.EXPECT().Create(ctx, orders_repo.CreateIn(serviceIn)).Return(primitive.ObjectID{}, assert.AnError)
			},
			wantErr: orders_service.ErrCannotCreateOrder,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockOrdersRepo := orders_repo_mocks.NewMockRepo(ctrl)
			tc.mockBehavior(mockOrdersRepo)

			s := orders_service.New(mockOrdersRepo)

			got, err := s.Create(ctx, serviceIn)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

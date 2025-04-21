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

func TestService_Find(t *testing.T) {
	var (
		ctx    = context.TODO()
		offset = 0
		limit  = 10
	)

	sampleOrder := orders_repo.OrderWithClientData{
		Title:          "Title",
		Description:    "Description",
		CompletionTime: 123456,
		Cost:           100,
		ClientName:     "Name",
		Rating:         4.7,
	}
	expectedOrder := orders_service.OrderWithClientData{
		Title:          sampleOrder.Title,
		Description:    sampleOrder.Description,
		CompletionTime: sampleOrder.CompletionTime,
		Cost:           sampleOrder.Cost,
		ClientName:     sampleOrder.ClientName,
		Rating:         sampleOrder.Rating,
	}

	type MockBehavior func(repo *orders_repo_mocks.MockRepo)

	tests := []struct {
		name         string
		mockBehavior MockBehavior
		want         orders_service.FindOut
		wantErr      error
	}{
		{
			name: "successful find",
			mockBehavior: func(r *orders_repo_mocks.MockRepo) {
				r.EXPECT().Find(ctx, offset, limit).Return(orders_repo.FindOut{Orders: []orders_repo.OrderWithClientData{sampleOrder}}, nil)
			},
			want: orders_service.FindOut{Orders: []orders_service.OrderWithClientData{expectedOrder}},
		},
		{
			name: "find failure",
			mockBehavior: func(r *orders_repo_mocks.MockRepo) {
				r.EXPECT().Find(ctx, offset, limit).Return(orders_repo.FindOut{}, assert.AnError)
			},
			want:    orders_service.FindOut{},
			wantErr: orders_service.ErrCannotFindOrders,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockRepo := orders_repo_mocks.NewMockRepo(ctrl)

			tc.mockBehavior(mockRepo)

			svc := orders_service.New(mockRepo)

			got, err := svc.Find(ctx, offset, limit)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	ctx := context.TODO()
	orderID := primitive.NewObjectID()

	sampleOrder := orders_repo.OrderWithClientData{
		Title:          "Title",
		Description:    "Description",
		CompletionTime: 123456,
		Cost:           100,
		ClientName:     "Name",
		Rating:         4.7,
	}
	expectedOrder := orders_service.OrderWithClientData(sampleOrder)

	type MockBehavior func(repo *orders_repo_mocks.MockRepo)

	tests := []struct {
		name         string
		mockBehavior MockBehavior
		want         orders_service.OrderWithClientData
		wantErr      error
	}{
		{
			name: "successful get",
			mockBehavior: func(r *orders_repo_mocks.MockRepo) {
				r.EXPECT().GetByID(ctx, orderID).Return(sampleOrder, nil)
			},
			want: expectedOrder,
		},
		{
			name: "order not found",
			mockBehavior: func(r *orders_repo_mocks.MockRepo) {
				r.EXPECT().GetByID(ctx, orderID).Return(orders_repo.OrderWithClientData{}, orders_repo.ErrOrderNotFound)
			},
			want:    orders_service.OrderWithClientData{},
			wantErr: orders_service.ErrOrderNotFound,
		},
		{
			name: "unexpected error",
			mockBehavior: func(r *orders_repo_mocks.MockRepo) {
				r.EXPECT().GetByID(ctx, orderID).Return(orders_repo.OrderWithClientData{}, assert.AnError)
			},
			want:    orders_service.OrderWithClientData{},
			wantErr: orders_service.ErrCannotGetOrder,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockRepo := orders_repo_mocks.NewMockRepo(ctrl)

			tc.mockBehavior(mockRepo)

			svc := orders_service.New(mockRepo)

			got, err := svc.GetByID(ctx, orderID)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

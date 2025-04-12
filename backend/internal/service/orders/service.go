package orders

import (
	"context"
	"errors"

	"github.com/moevm/nosql1h25-writer/backend/internal/repo/orders"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	ordersRepo orders.Repo
}

func New(ordersRepo orders.Repo) Service {
	return &service{ordersRepo: ordersRepo}
}

func (s *service) FindOrders(ctx context.Context, offset, limit int64) ([]orders.FindOrdersOut, int64, error) {
	if limit <= 0 || limit > 100 || offset < 0 {
		return nil, 0, ErrInvalidPagination
	}

	orders, total, err := s.ordersRepo.FindOrders(ctx, offset, limit)
	if err != nil {
		if err == ErrOrdersNotFound {
			return nil, 0, ErrOrdersNotFound
		}
		return nil, 0, ErrCannotGetOrders
	}

	return orders, total, nil
}

func (s *service) GetByID(ctx context.Context, id string) (orders.FindOrdersOut, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return orders.FindOrdersOut{}, ErrOrdersNotFound
	}
	order, err := s.ordersRepo.GetByID(ctx, objID)
	if err != nil {
		if errors.Is(err, ErrOrdersNotFound) {
			return orders.FindOrdersOut{}, ErrOrdersNotFound
		}
		return orders.FindOrdersOut{}, ErrCannotGetOrders
	}
	return order, nil
}

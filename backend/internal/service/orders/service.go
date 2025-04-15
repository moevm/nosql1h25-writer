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

func (s *service) Find(ctx context.Context, offset, limit int) (FindOut, error) {
	repoFindOut, err := s.ordersRepo.Find(ctx, offset, limit)
	if err != nil {
		return FindOut{}, err
	}
	var serviceFindOut FindOut
	for _, order := range repoFindOut.Orders {
		serviceFindOut.Orders = append(serviceFindOut.Orders, OrderWithClientData{
			Title:          order.Title,
			Description:    order.Description,
			CompletionTime: order.CompletionTime,
			Cost:           order.Cost,
			ClientName:     order.ClientName,
			Rating:         order.Rating,
		})
	}
	return serviceFindOut, nil
}

func (s *service) GetByID(ctx context.Context, id primitive.ObjectID) (OrderWithClientData, error) {
	getByIDOut, err := s.ordersRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return OrderWithClientData{}, ErrOrderNotFound
		}
		return OrderWithClientData{}, ErrCannotGetOrder
	}
	return OrderWithClientData(getByIDOut), nil
}

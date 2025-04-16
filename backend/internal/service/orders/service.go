package orders

import (
	"context"
	"errors"

	"github.com/moevm/nosql1h25-writer/backend/internal/repo/orders"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	ordersRepo orders.Repo
}

func New(ordersRepo orders.Repo) Service {
	return &service{ordersRepo: ordersRepo}
}

func (s *service) Find(ctx context.Context, offset, limit int) (FindOut, error) {
	out, err := s.ordersRepo.Find(ctx, offset, limit)
	if err != nil {
		log.Errorf("OrderService.Find - s.ordersRepo: %v", err)
		return FindOut{}, ErrCannotFindOrders
	}
	var serviceFindOut FindOut
	for _, order := range out.Orders {
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
	out, err := s.ordersRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return OrderWithClientData{}, ErrOrderNotFound
		}
		log.Errorf("OrderService.Get - s.ordersRepo: %v", err)
		return OrderWithClientData{}, ErrCannotGetOrder
	}
	return OrderWithClientData(out), nil
}

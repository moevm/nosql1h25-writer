package orders

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/repo/orders"
)

type service struct {
	ordersRepo orders.Repo
}

func New(ordersRepo orders.Repo) Service {
	return &service{ordersRepo: ordersRepo}
}

func (s *service) Find(ctx context.Context, offset, limit int, minCost, maxCost *int, sortBy *string) (FindOut, error) {
	out, err := s.ordersRepo.Find(ctx, offset, limit, minCost, maxCost, sortBy)
	if err != nil {
		log.Errorf("OrderService.Find - s.ordersRepo: %v", err)
		return FindOut{}, ErrCannotFindOrders
	}
	var serviceFindOut FindOut
	serviceFindOut.Total = out.Total
	for _, order := range out.Orders {
		serviceFindOut.Orders = append(serviceFindOut.Orders, OrderWithClientData{
			ID:             order.ID,
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

func (s *service) Create(ctx context.Context, in CreateIn) (primitive.ObjectID, error) {
	id, err := s.ordersRepo.Create(ctx, orders.CreateIn{
		ClientID:       in.ClientID,
		Title:          in.Title,
		Description:    in.Description,
		CompletionTime: in.CompletionTime,
		Cost:           in.Cost,
	})
	if err != nil {
		log.Errorf("service.orders.Create - s.ordersRepo.Create: %v", err)
		return primitive.ObjectID{}, ErrCannotCreateOrder
	}

	return id, nil
}

func (s *service) Update(ctx context.Context, in UpdateIn) error {
	err := s.ordersRepo.Update(ctx, orders.UpdateIn{
		OrderID:        in.OrderID,
		Title:          in.Title,
		Description:    in.Description,
		CompletionTime: in.CompletionTime,
		Cost:           in.Cost,
	})
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return ErrOrderNotFound
		}
		log.Errorf("service.orders.Update - s.ordersRepo.Update: %v", err)
		return ErrCannotUpdateOrder
	}

	return nil
}

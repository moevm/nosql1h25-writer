package orders

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/repo/orders"
	userservice "github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type orderService struct {
	ordersRepo   orders.Repo
	usersService userservice.Service
}

func New(ordersRepo orders.Repo, usersService userservice.Service) Service {
	return &orderService{
		ordersRepo:   ordersRepo,
		usersService: usersService,
	}
}

func (s *orderService) Find(ctx context.Context, offset, limit int, minCost, maxCost *int, sortBy *string) (FindOut, error) {
	out, err := s.ordersRepo.Find(ctx, offset, limit, minCost, maxCost, sortBy)
	if err != nil {
		logrus.Errorf("OrderService.Find - s.ordersRepo.Find: %v", err)
		return FindOut{}, ErrCannotFindOrders
	}

	serviceFindOut := FindOut{Total: out.Total}
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

func (s *orderService) GetByID(ctx context.Context, id primitive.ObjectID) (OrderWithClientData, error) {
	out, err := s.ordersRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return OrderWithClientData{}, ErrOrderNotFound
		}
		logrus.Errorf("OrderService.GetByID - s.ordersRepo.GetByID: %v", err)
		return OrderWithClientData{}, ErrCannotGetOrder
	}
	return OrderWithClientData(out), nil
}

func (s *orderService) GetByIDExt(ctx context.Context, id primitive.ObjectID) (entity.OrderExt, error) {
	orderExt, err := s.ordersRepo.GetOrderExtByID(ctx, id)
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return entity.OrderExt{}, ErrOrderNotFound
		}
		logrus.Errorf("OrderService.GetByIDExt - s.ordersRepo.GetOrderExtByID: %v", err)
		return entity.OrderExt{}, ErrCannotGetOrder
	}
	return orderExt, nil
}

func (s *orderService) Create(ctx context.Context, in CreateIn) (primitive.ObjectID, error) {
	id, err := s.ordersRepo.Create(ctx, orders.CreateIn{
		ClientID:       in.ClientID,
		Title:          in.Title,
		Description:    in.Description,
		CompletionTime: in.CompletionTime,
		Cost:           in.Cost,
	})
	if err != nil {
		logrus.Errorf("OrderService.Create - s.ordersRepo.Create: %v", err)
		return primitive.ObjectID{}, ErrCannotCreateOrder
	}
	return id, nil
}

func (s *orderService) Response(ctx context.Context, orderID, userID primitive.ObjectID) (primitive.ObjectID, error) {
	orderExt, err := s.GetByIDExt(ctx, orderID)
	if err != nil {
		logrus.Errorf("OrderService.Response - GetByIDExt: %v", err)
		return primitive.ObjectID{}, err
	}

	userExt, err := s.usersService.GetByIDExt(ctx, userID)
	if err != nil {
		logrus.Errorf("OrderService.Response - usersService.GetByIDExt: %v", err)
		return primitive.ObjectID{}, err
	}

	response := entity.Response{
		FreelancerID:   userID,
		FreelancerName: userExt.User.DisplayName,
		ChatID:         primitive.NewObjectID(),
		CoverLetter:    "I would like to work on this project", // TODO: параметризовать
		Active:         true,
		CreatedAt:      time.Now().UTC(),
	}

	orderExt.Responses = append(orderExt.Responses, response)
	
	if err := s.ordersRepo.Update(ctx, orderID, orderExt); err != nil {
		logrus.Errorf("OrderService.Response - s.ordersRepo.Update: %v", err)
		return primitive.ObjectID{}, ErrCannotUpdateOrder
	}

	return orderID, nil
}
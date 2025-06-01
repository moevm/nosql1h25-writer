package orders

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/repo/orders"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type service struct {
	ordersRepo   orders.Repo
	usersService users.Service
}

func New(ordersRepo orders.Repo, usersService users.Service) Service {
	return &service{
		ordersRepo:   ordersRepo,
		usersService: usersService,
	}
}

func (s *service) Find(ctx context.Context, in FindIn) (FindOut, error) {
	out, err := s.ordersRepo.Find(ctx, orders.FindIn{
		Limit:   in.Limit,
		Offset:  in.Offset,
		MinCost: in.MinCost,
		MaxCost: in.MaxCost,
		MinTime: in.MinTime,
		MaxTime: in.MaxTime,
		Search:  in.Search,
		SortBy:  in.SortBy,
	})
	if err != nil {
		log.Errorf("orders.service.Find - s.ordersRepo.Find: %v", err)
		return FindOut{}, ErrCannotFindOrders
	}

	return FindOut{Total: out.Total, Orders: out.Orders}, nil
}

func (s *service) GetByID(ctx context.Context, id primitive.ObjectID) (OrderWithClientData, error) {
	out, err := s.ordersRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return OrderWithClientData{}, ErrOrderNotFound
		}

		log.Errorf("orders.service.GetByID - s.ordersRepo.GetByID: %v", err)
		return OrderWithClientData{}, ErrCannotGetOrder
	}

	return OrderWithClientData(out), nil
}

func (s *service) GetByIDExt(ctx context.Context, id primitive.ObjectID) (entity.OrderExt, error) {
	out, err := s.ordersRepo.GetByIDExt(ctx, id)
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return entity.OrderExt{}, ErrOrderNotFound
		}

		log.Errorf("orders.service.GetByIDExt - s.ordersRepo.GetByIDExt: %v", err)
		return entity.OrderExt{}, ErrCannotGetOrder
	}

	return out, nil
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
		Status:         in.Status,
		FreelancerID:   in.FreelancerID,
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

func (s *service) CreateResponse(ctx context.Context, orderID primitive.ObjectID, userID primitive.ObjectID, coverLetter string) error {
	orderExt, err := s.GetByIDExt(ctx, orderID)
	if err != nil {
		return err
	}

	if orderExt.ClientID == userID || orderExt.Statuses[len(orderExt.Statuses)-1].Type != entity.StatusTypeBeginning {
		return ErrCannotResponse
	}

	for _, response := range orderExt.Responses {
		if response.Active && response.FreelancerID == userID {
			return ErrCannotResponse
		}
	}

	user, err := s.usersService.GetByIDExt(ctx, userID)
	if err != nil {
		return err
	}

	err = s.ordersRepo.CreateResponse(ctx, orderID, userID, coverLetter, user.DisplayName)
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return ErrOrderNotFound
		}

		log.Errorf("orders.service.CreateResponse - s.ordersRepo.CreateResponse: %v", err)
		return ErrCannotCreateResponse
	}

	return nil
}

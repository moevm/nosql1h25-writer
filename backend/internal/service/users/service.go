package users

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/repo/orders"
	"github.com/moevm/nosql1h25-writer/backend/internal/repo/users"
)

type service struct {
	usersRepo  users.Repo
	ordersRepo orders.Repo
}

func New(usersRepo users.Repo, ordersRepo orders.Repo) Service {
	return &service{
		usersRepo:  usersRepo,
		ordersRepo: ordersRepo,
	}
}

func (s *service) Find(ctx context.Context, in FindIn) (FindOut, error) {
	out, err := s.usersRepo.Find(ctx, users.FindIn{
		Limit:               in.Limit,
		Offset:              in.Offset,
		NameSearch:          in.NameSearch,
		EmailSearch:         in.EmailSearch,
		Roles:               in.Roles,
		MinFreelancerRating: in.MinFreelancerRating,
		MaxFreelancerRating: in.MaxFreelancerRating,
		MinClientRating:     in.MinClientRating,
		MaxClientRating:     in.MaxClientRating,
		MinCreatedAt:        in.MinCreatedAt,
		MaxCreatedAt:        in.MaxCreatedAt,
		MinBalance:          in.MinBalance,
		MaxBalance:          in.MaxBalance,
		SortBy:              in.SortBy,
	})
	if err != nil {
		log.Errorf("users.service.Find - s.usersRepo.Find: %v", err)
		return FindOut{}, ErrCannotFindUsers
	}

	return FindOut{Users: out.Users, Total: out.Total}, nil
}

func (s *service) UpdateBalance(ctx context.Context, userID primitive.ObjectID, op OperationType, amount int) (int, error) {
	var (
		newBalance int
		err        error
	)

	switch op {
	case OperationTypeDeposit:
		newBalance, err = s.usersRepo.Deposit(ctx, userID, amount)
	case OperationTypeWithdraw:
		newBalance, err = s.usersRepo.Withdraw(ctx, userID, amount)
	}

	if err != nil {
		if errors.Is(err, users.ErrInsufficientFunds) {
			return 0, ErrInsufficientFunds
		}

		log.Errorf("Service.UpdateBalance - s.usersRepo.%s: %v", op, err)
		return 0, ErrUpdateBalance
	}

	return newBalance, nil
}

func (s *service) GetByIDExt(ctx context.Context, userID primitive.ObjectID) (entity.UserExt, error) {
	user, err := s.usersRepo.GetByIDExt(ctx, userID)

	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return entity.UserExt{}, ErrUserNotFound
		}

		log.Errorf("users.service.GetByIDExt - s.usersRepo.GetByIDExt: %v", err)
		return entity.UserExt{}, ErrCannotGetUser
	}
	return user, nil
}

func (s *service) Update(ctx context.Context, in UpdateIn) error {
	err := s.usersRepo.Update(ctx, users.UpdateIn{
		UserID:                in.UserID,
		DisplayName:           in.DisplayName,
		FreelancerDescription: in.FreelancerDescription,
		ClientDescription:     in.ClientDescription,
	})
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return ErrUserNotFound
		}

		log.Errorf("users.service.Update - usersRepo.Update: %v", err)
		return ErrCannotUpdateUser
	}

	return nil
}

func (s *service) FindOrdersByUserID(ctx context.Context, requesterID, targetUserID primitive.ObjectID) ([]entity.OrderExt, error) {
	orders, err := s.ordersRepo.FindByUserIDExt(ctx, targetUserID)
	if err != nil {
		log.Errorf("users.service.FindOrdersByUserID - s.ordersRepo.FindByUserIDExt: %v", err)
		return nil, ErrCannotFindOrders
	}

	return orders, nil
}

func (s *service) FindOrdersByResponseUserID(ctx context.Context, freelancerID primitive.ObjectID) ([]entity.OrderExt, error) {
	orders, err := s.ordersRepo.FindByResponseUserID(ctx, freelancerID)
	if err != nil {
		log.Errorf("users.service.FindOrdersByResponseUserID - s.ordersRepo.FindByResponseUserID: %v", err)
		return nil, ErrCannotFindOrders
	}

	return orders, nil
}

package users

import (
	"context"
	"errors"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/repo/users"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	usersRepo users.Repo
}

func New(repo users.Repo) Service {
	return &service{usersRepo: repo}
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

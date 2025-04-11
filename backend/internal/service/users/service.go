package users

import (
	"context"
	"errors"

	"github.com/moevm/nosql1h25-writer/backend/internal/repo/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
	log "github.com/sirupsen/logrus"
)

type service struct {
	repo users.Repo
}

func New(repo users.Repo) Service {
	return &service{repo: repo}
}

func (s *service) UpdateBalance(ctx context.Context, userID primitive.ObjectID, op OperationType, amount int) error {
	switch op {
	case OperationTypeDeposit:
		if err := s.repo.Deposit(ctx, userID, amount); err != nil {
			if errors.Is(err, users.ErrCannotDeposit) {
				return ErrCannotDeposit
			}
			log.Errorf("Service.UpdateBalance (deposit) - s.repo.Deposit: %v", err)
			return ErrCannotDeposit
		}
	case OperationTypeWithdraw:
		if err := s.repo.Withdraw(ctx, userID, amount); err != nil {
			if errors.Is(err, users.ErrInsufficientFunds) {
				return ErrInsufficientFunds
			}
			log.Errorf("Service.UpdateBalance (withdraw) - s.repo.Withdraw: %v", err)
			return ErrCannotWithdraw
		}
	default:
		return errors.New("invalid operation type")
	}
	return nil
}

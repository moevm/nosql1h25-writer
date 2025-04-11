package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	UpdateBalance(ctx context.Context, userID primitive.ObjectID, op OperationType, amount int) (int, error)
}

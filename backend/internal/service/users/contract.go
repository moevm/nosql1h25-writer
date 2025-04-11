package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	UpdateBalance(ctx context.Context, userID primitive.ObjectID, op OperationType, amount int) error
}

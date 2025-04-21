package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Repo
type Repo interface {
	Create(ctx context.Context, in CreateIn) (primitive.ObjectID, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (entity.User, error)
	Deposit(ctx context.Context, userID primitive.ObjectID, amount int) (int, error)
	Withdraw(ctx context.Context, userID primitive.ObjectID, amount int) (int, error)
}

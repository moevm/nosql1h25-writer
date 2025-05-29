package users

import (
	"context"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	UpdateBalance(ctx context.Context, userID primitive.ObjectID, op OperationType, amount int) (int, error)
	GetByIDExt(ctx context.Context, userID primitive.ObjectID) (entity.UserExt, error)
	Update(ctx context.Context, in UpdateIn) error
}

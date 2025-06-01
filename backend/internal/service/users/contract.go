package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	Find(ctx context.Context, in FindIn) (FindOut, error)
	UpdateBalance(ctx context.Context, userID primitive.ObjectID, op OperationType, amount int) (int, error)
	GetByIDExt(ctx context.Context, userID primitive.ObjectID) (entity.UserExt, error)
	Update(ctx context.Context, in UpdateIn) error
	FindOrdersByUserID(ctx context.Context, requesterID, targetUserID primitive.ObjectID) ([]entity.OrderExt, error)
	FindOrdersByResponseUserID(ctx context.Context, freelancerID primitive.ObjectID) ([]entity.OrderExt, error)
}

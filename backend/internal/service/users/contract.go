package users

import (
	"context"

<<<<<<< HEAD
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type Service interface {
	FindUsers(ctx context.Context, params entity.UserSearchParams) ([]entity.UserExt, int64, error)
=======
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	UpdateBalance(ctx context.Context, userID primitive.ObjectID, op OperationType, amount int) (int, error)
>>>>>>> origin/main
}

package usersExt

import (
	"context"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	FindUsers(ctx context.Context, params entity.UserSearchParams) ([]entity.UserExt, int64, error)
	FindUserByID(ctx context.Context, userID primitive.ObjectID) (entity.UserExt, error)
}

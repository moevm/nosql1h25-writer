package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type Repo interface {
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetById(ctx context.Context, ID primitive.ObjectID) (entity.User, error)
}

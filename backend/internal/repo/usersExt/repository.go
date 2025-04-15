package usersExt

import (
	"context"
	"errors"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrUserNotFound = errors.New("user not found")

type Repo interface {
	FindUsers(ctx context.Context, params entity.UserSearchParams) ([]entity.UserExt, int64, error)

	FindUserByID(ctx context.Context, userID primitive.ObjectID) (entity.UserExt, error)
}

func New(db mongoifc.Database) Repo {
	return &mongodbRepo{db: db}
}

type mongodbRepo struct {
	db mongoifc.Database
}

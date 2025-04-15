package usersExt

import (
	"context"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"

	"github.com/sv-tools/mongoifc"
)

type Repo interface {
	FindUsers(ctx context.Context, params entity.UserSearchParams) ([]entity.UserExt, int64, error)
}

func New(db mongoifc.Database) Repo {
	return &mongodbRepo{db: db}
}

type mongodbRepo struct {
	db mongoifc.Database
}

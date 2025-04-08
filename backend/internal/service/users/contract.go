package users

import (
	"context"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type Service interface {
	FindUsers(ctx context.Context, params entity.UserSearchParams) ([]entity.UserExt, int64, error)
}

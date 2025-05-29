package orders

import (
	"context"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Repo
type Repo interface {
	Create(ctx context.Context, in CreateIn) (primitive.ObjectID, error)
	Find(ctx context.Context, offset, limit int, minCost, maxCost *int, sortBy *string) (FindOut, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (OrderWithClientData, error)
	GetOrderExtByID(ctx context.Context, id primitive.ObjectID) (entity.OrderExt, error)
	Update(ctx context.Context, id primitive.ObjectID, order entity.OrderExt) error
}

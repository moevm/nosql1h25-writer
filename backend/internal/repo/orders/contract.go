package orders

import (
	"context"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Repo
type Repo interface {
	// FindOrders(ctx context.Context, offset, limit int64) ([]FindOrdersOut, int64, error)
	// GetByID(ctx context.Context, id primitive.ObjectID) (FindOrdersOut, error)
	Find(ctx context.Context, offset, limit int) (FindOut, error)
}

package orders

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Repo
type Repo interface {
	Find(ctx context.Context, offset, limit int) (FindOut, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (OrderWithClientData, error)
}

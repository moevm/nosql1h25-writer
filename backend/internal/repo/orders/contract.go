package orders

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Repo
type Repo interface {
	Create(ctx context.Context, in CreateIn) (primitive.ObjectID, error)
}

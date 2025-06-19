package stats

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Repo
type Repo interface {
	Aggregate(ctx context.Context, x, y string) ([]primitive.ObjectID, error)
}

package orders

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	Create(ctx context.Context, in CreateIn) (primitive.ObjectID, error)
	Find(ctx context.Context, offset, limit int, minCost, maxCost *int, sortBy *string) (FindOut, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (OrderWithClientData, error)
	Update(ctx context.Context, in UpdateIn) error
}

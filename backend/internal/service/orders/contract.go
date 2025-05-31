package orders

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	Create(ctx context.Context, in CreateIn) (primitive.ObjectID, error)
	Find(ctx context.Context, offset, limit int, minCost, maxCost *int, sortBy *string) (FindOut, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (OrderWithClientData, error)
	GetByIDExt(ctx context.Context, id primitive.ObjectID) (entity.OrderExt, error)
	Update(ctx context.Context, in UpdateIn) error
	CreateResponse(ctx context.Context, orderID primitive.ObjectID, userID primitive.ObjectID, coverletter string) error
}

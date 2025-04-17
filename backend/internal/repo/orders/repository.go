package orders

import (
	"context"

	"github.com/jonboulle/clockwork"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type repository struct {
	ordersColl mongoifc.Collection
	clock      clockwork.Clock
}

func New(ordersColl mongoifc.Collection, clock clockwork.Clock) Repo {
	return &repository{
		ordersColl: ordersColl,
		clock:      clock,
	}
}

func (r *repository) Create(ctx context.Context, in CreateIn) (primitive.ObjectID, error) {
	now := r.clock.Now()
	order := entity.DefaultOrder(in.ClientID, in.Title, in.Description, in.CompletionTime, in.Cost, now, now)

	res, err := r.ordersColl.InsertOne(ctx, order)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return res.InsertedID.(primitive.ObjectID), nil //nolint:forcetypeassert
}

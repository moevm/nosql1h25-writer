package stats

import (
	"context"

	"github.com/jonboulle/clockwork"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type repository struct {
	usersColl  mongoifc.Collection
	ordersColl mongoifc.Collection
	clock      clockwork.Clock
}

func New(usersColl mongoifc.Collection, ordersColl mongoifc.Collection, clock clockwork.Clock) Repo {
	return &repository{
		usersColl:  usersColl,
		ordersColl: ordersColl,
		clock:      clock,
	}
}

func (r *repository) Aggregate(ctx context.Context, x, y string) ([]primitive.ObjectID, error) {
	return nil, nil
}

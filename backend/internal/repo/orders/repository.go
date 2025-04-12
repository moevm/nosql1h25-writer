package orders

import (
	"context"
	"errors"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	ordersColl mongoifc.Collection
}

func New(ordersColl mongoifc.Collection) Repo {
	return &repository{ordersColl: ordersColl}
}

func (r *repository) FindOrders(ctx context.Context, offset, limit int64) ([]FindOrdersOut, int64, error) {
	filter := bson.D{}
	findOpts := options.Find().
		SetSkip(offset).
		SetLimit(limit)

	cursor, err := r.ordersColl.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, ErrCannotGetOrders
	}
	defer cursor.Close(ctx)

	var orders []FindOrdersOut
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, 0, ErrCannotGetOrders
	}

	total, err := r.ordersColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, ErrCannotGetOrders
	}

	if len(orders) == 0 {
		return nil, 0, ErrNotFound
	}

	return orders, total, nil
}

func (r *repository) GetByID(ctx context.Context, id primitive.ObjectID) (FindOrdersOut, error) {
	var order FindOrdersOut
	err := r.ordersColl.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return FindOrdersOut{}, ErrNotFound
		}
		return FindOrdersOut{}, ErrCannotGetOrders
	}
	return order, nil
}

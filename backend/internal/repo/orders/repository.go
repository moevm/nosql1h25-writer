package orders

import (
	"context"
	"errors"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	ordersColl mongoifc.Collection
}

func New(ordersColl mongoifc.Collection) Repo {
	return &repository{ordersColl: ordersColl}
}

func (r *repository) Find(ctx context.Context, offset, limit int) (FindOut, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"active": true,
		}}},
		{{Key: "$match", Value: bson.M{
			"$expr": bson.M{
				"$eq": bson.A{
					bson.M{"$arrayElemAt": bson.A{"$statuses.title", -1}},
					entity.StatusTypeBeginning,
				},
			},
		}}},
		{{Key: "$skip", Value: int64(offset)}},
		{{Key: "$limit", Value: int64(limit)}},
	}
	cursor, err := r.ordersColl.Aggregate(ctx, pipeline)
	if err != nil {
		return FindOut{}, err
	}

	var orders []entity.OrderExt
	if err := cursor.All(ctx, &orders); err != nil {
		return FindOut{}, err
	}

	var out FindOut
	for _, order := range orders {
		out.Orders = append(out.Orders, OrderWithClientData{
			Title:          order.Title,
			Description:    order.Description,
			CompletionTime: int(order.CompletionTime),
			Cost:           order.Cost,
			ClientName:     "",  // пока без объединения с юзерами
			Rating:         0.0, // пока без объединения с юзерами
		})
	}
	return out, nil
}

func (r *repository) GetByID(ctx context.Context, id primitive.ObjectID) (OrderWithClientData, error) {
	var order entity.Order
	err := r.ordersColl.FindOne(ctx, bson.M{"_id": id, "active": true}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return OrderWithClientData{}, ErrOrderNotFound
		}
		return OrderWithClientData{}, err
	}
	return OrderWithClientData{
		Title:          order.Title,
		Description:    order.Description,
		CompletionTime: int(order.CompletionTime),
		Cost:           order.Cost,
		ClientName:     "",  // пока без объединения с юзерами
		Rating:         0.0, // пока без объединения с юзерами
	}, nil
}

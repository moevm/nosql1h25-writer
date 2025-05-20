package orders

import (
	"context"
	"errors"

	"github.com/jonboulle/clockwork"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

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

func (r *repository) Find(ctx context.Context, offset, limit int, minCost, maxCost *int) (FindOut, error) {
	matchFilter := bson.M{
		"active": true,
	}

	costCond := bson.M{}
	if minCost != nil {
		costCond["$gte"] = *minCost
	}
	if maxCost != nil {
		costCond["$lte"] = *maxCost
	}
	if len(costCond) > 0 {
		matchFilter["cost"] = costCond
	}

	// pipeline := mongo.Pipeline{
	// 	{{Key: "$match", Value: bson.M{
	// 		"active": true,
	// 	}}},
	// 	{{Key: "$match", Value: bson.M{
	// 		"$expr": bson.M{
	// 			"$eq": bson.A{
	// 				bson.M{"$arrayElemAt": bson.A{"$statuses.type", -1}},
	// 				entity.StatusTypeBeginning,
	// 			},
	// 		},
	// 	}}},
	// 	{{Key: "$facet", Value: bson.M{
	// 		"orders": mongo.Pipeline{
	// 			{{Key: "$project", Value: bson.M{
	// 				"_id":            1,
	// 				"clientId":       1,
	// 				"title":          1,
	// 				"description":    1,
	// 				"completionTime": 1,
	// 				"cost":           1,
	// 				"freelancerId":   1,
	// 				"budget":         1,
	// 				"createdAt":      1,
	// 				"updatedAt":      1,
	// 				"responses":      1,
	// 				"statuses":       1,
	// 			}}},
	// 			{{Key: "$skip", Value: int64(offset)}},
	// 			{{Key: "$limit", Value: int64(limit)}},
	// 		},
	// 		"total": mongo.Pipeline{
	// 			{{Key: "$count", Value: "count"}},
	// 		},
	// 	}}},
	// }

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchFilter}},
		{{Key: "$match", Value: bson.M{
			"$expr": bson.M{
				"$eq": bson.A{
					bson.M{"$arrayElemAt": bson.A{"$statuses.type", -1}},
					entity.StatusTypeBeginning,
				},
			},
		}}},
		{{Key: "$facet", Value: bson.M{
			"orders": mongo.Pipeline{
				{{Key: "$project", Value: bson.M{
					"_id":            1,
					"clientId":       1,
					"title":          1,
					"description":    1,
					"completionTime": 1,
					"cost":           1,
					"freelancerId":   1,
					"budget":         1,
					"createdAt":      1,
					"updatedAt":      1,
					"responses":      1,
					"statuses":       1,
				}}},
				{{Key: "$skip", Value: int64(offset)}},
				{{Key: "$limit", Value: int64(limit)}},
			},
			"total": mongo.Pipeline{
				{{Key: "$count", Value: "count"}},
			},
		}}},
	}

	cursor, err := r.ordersColl.Aggregate(ctx, pipeline)
	if err != nil {
		return FindOut{}, err
	}

	var result []struct {
		Orders []entity.OrderExt `bson:"orders"`
		Total  []struct {
			Count int `bson:"count"`
		} `bson:"total"`
	}

	if err := cursor.All(ctx, &result); err != nil {
		return FindOut{}, err
	}

	totalCount := 0
	if len(result) > 0 && len(result[0].Total) > 0 {
		totalCount = result[0].Total[0].Count
	}

	var out FindOut
	out.Total = totalCount
	for _, order := range result[0].Orders {
		out.Orders = append(out.Orders, OrderWithClientData{
			ID:             order.ID.Hex(),
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

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

func (r *repository) Find(ctx context.Context, offset, limit int, minCost, maxCost *int, sortBy *string) (FindOut, error) {
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

	var sortStage bson.D
	if sortBy == nil {
		sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}}
	} else {
		switch *sortBy {
		case "newest":
			sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}}
		case "oldest":
			sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: 1}}}}
		case "cost_asc":
			sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: "cost", Value: 1}}}}
		case "cost_desc":
			sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: "cost", Value: -1}}}}
		default:
			sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}} // fallback
		}
	}

	ordersPipeline := mongo.Pipeline{
		{{Key: "$project", Value: bson.M{
			"_id":            1,
			"clientId":       1,
			"title":          1,
			"description":    1,
			"completionTime": 1,
			"cost":           1,
			"freelancerId":   1,
			"createdAt":      1,
			"updatedAt":      1,
			"responses":      1,
			"statuses":       1,
		}}},
	}

	if len(sortStage) > 0 {
		ordersPipeline = append(ordersPipeline, sortStage)
	}

	ordersPipeline = append(ordersPipeline,
		bson.D{{Key: "$skip", Value: int64(offset)}},
		bson.D{{Key: "$limit", Value: int64(limit)}},
	)

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
			"orders": ordersPipeline,
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
			ID:             order.ID,
			Title:          order.Title,
			Description:    order.Description,
			CompletionTime: order.CompletionTime,
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
		ClientID:       order.ClientID,
		Title:          order.Title,
		Description:    order.Description,
		CompletionTime: order.CompletionTime,
		Cost:           order.Cost,
		ClientName:     "",  // пока без объединения с юзерами
		Rating:         0.0, // пока без объединения с юзерами
	}, nil
}

func (r *repository) GetByIDExt(ctx context.Context, id primitive.ObjectID) (order entity.OrderExt, _ error) {
	if err := r.ordersColl.FindOne(ctx, bson.M{"_id": id, "active": true}).Decode(&order); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return order, ErrOrderNotFound
		}
		return order, err
	}
	return order, nil
}

func (r *repository) CreateResponse(
	ctx context.Context,
	orderID primitive.ObjectID,
	userID primitive.ObjectID,
	coverLetter, freelancerName string,
) error {
	now := r.clock.Now()

	response := entity.Response{
		FreelancerName: freelancerName,
		CoverLetter:    coverLetter,
		FreelancerID:   userID,
		CreatedAt:      now,
		Active:         true,
	}

	update := bson.M{
		"$push": bson.M{
			"responses": response,
		},
		"$set": bson.M{
			"updatedAt": now,
		},
	}

	_, err := r.ordersColl.UpdateOne(
		ctx,
		bson.M{"_id": orderID, "active": true},
		update,
	)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrOrderNotFound
		}
		return err
	}

	return nil
}

func (r *repository) Update(ctx context.Context, in UpdateIn) error {
	now := r.clock.Now()

	updateSet := bson.M{
		"updatedAt": now,
	}

	if in.Title != nil {
		updateSet["title"] = *in.Title
	}

	if in.Description != nil {
		updateSet["description"] = *in.Description
	}

	if in.CompletionTime != nil {
		updateSet["completionTime"] = *in.CompletionTime
	}

	if in.Cost != nil {
		updateSet["cost"] = *in.Cost
	}

	if in.FreelancerID != nil {
		updateSet["freelancerId"] = *in.FreelancerID
	}

	update := bson.M{
		"$set": updateSet,
	}

	// если статус задан, пушим его в массив statuses
	if in.Status != nil {
		update["$push"] = bson.M{
			"statuses": entity.Status{
				Type:      *in.Status,
				CreatedAt: now,
			},
		}
	}

	err := r.ordersColl.FindOneAndUpdate(
		ctx,
		bson.M{"_id": in.OrderID, "active": true},
		update,
	).Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrOrderNotFound
		}
		return err
	}

	return nil
}

func (r *repository) FindByUserIDExt(ctx context.Context, userID primitive.ObjectID) ([]entity.OrderExt, error) {
	filter := bson.M{
		"active":   true,
		"clientId": userID,
	}

	cursor, err := r.ordersColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var orders []entity.OrderExt
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

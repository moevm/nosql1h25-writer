package orders

import (
	"context"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	ordersColl mongoifc.Collection
}

func New(ordersColl mongoifc.Collection) Repo {
	return &repository{ordersColl: ordersColl}
}

func (r *repository) Find(ctx context.Context, offset, limit int) (FindOut, error) {
	filter := bson.M{
		"active": true,
	}

	findOpts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit))
	cursor, err := r.ordersColl.Find(ctx, filter, findOpts)
	if err != nil {
		return FindOut{}, err
	}
	defer cursor.Close(ctx)

	var ordersExt []entity.OrderExt
	for cursor.Next(ctx) {
		var orderExt entity.OrderExt
		if err := cursor.Decode(&orderExt); err != nil {
			return FindOut{}, err
		}
		ordersExt = append(ordersExt, orderExt)
	}

	var dto FindOut
	for _, order := range ordersExt {
		if len(order.Statuses) == 0 {
			continue
		}
		lastStatus := order.Statuses[len(order.Statuses)-1]
		if lastStatus.Type != entity.StatusTypeBeginning {
			continue
		}
		dto.Orders = append(dto.Orders, OrderWithClientData{
			Title:          order.Title,
			Description:    order.Description,
			CompletionTime: int(order.CompletionTime),
			Cost:           order.Cost,
			ClientName:     "",  //пока без объединения с юзерами
			Rating:         0.0, //пока без объединения с юзерами
		})
	}
	return dto, nil
}

// func (r *repository) GetByID(ctx context.Context, id primitive.ObjectID) (FindOrdersOut, error) {
// 	var order FindOrdersOut
// 	err := r.ordersColl.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return FindOrdersOut{}, ErrNotFound
// 		}
// 		return FindOrdersOut{}, ErrCannotGetOrders
// 	}
// 	return order, nil
// }

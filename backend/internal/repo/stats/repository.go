package stats

import (
	"context"
	"fmt"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type repository struct {
	usersColl  mongoifc.Collection
	ordersColl mongoifc.Collection
}

func New(usersColl mongoifc.Collection, ordersColl mongoifc.Collection) Repo {
	return &repository{
		usersColl:  usersColl,
		ordersColl: ordersColl,
	}
}

func buildPipeline(x, y string, aggType entity.Aggregation) (mongo.Pipeline, string, error) {
	infoX, okX := fieldInfo[x]
	if !okX {
		return nil, "", fmt.Errorf("invalid X field %s", x)
	}
	// определяем коллекцию: берём из X
	collName := infoX.Coll

	var pipe mongo.Pipeline

	// $project: задаём поля x и y
	projFields := bson.D{{Key: "x", Value: "$" + infoX.Path}}
	// для y: обрабатываем count или IsArray
	if y == "count" {
		// count без привязки к Path
		projFields = append(projFields, bson.E{Key: "y_count", Value: bson.D{{Key: "$literal", Value: 1}}})
	} else {
		infoY, okY := fieldInfo[y]
		if !okY {
			return nil, "", fmt.Errorf("invalid y field %s", y)
		}
		// проверка на совпадение коллекций
		if infoY.Coll != collName {
			return nil, "", fmt.Errorf("x and y are from different collections: %s vs %s", collName, infoY.Coll)
		}
		if infoY.IsArray {
			projFields = append(projFields, bson.E{Key: "y", Value: bson.D{
				{Key: "$size", Value: bson.D{
					{Key: "$ifNull", Value: bson.A{"$" + infoY.Path, bson.A{}}},
				}},
			}})
		} else {
			projFields = append(projFields, bson.E{Key: "y", Value: "$" + infoY.Path})
		}
	}
	pipe = append(pipe, bson.D{{Key: "$project", Value: projFields}})

	// $group: группируем по x
	var aggExpr primitive.D
	switch aggType {
	case entity.AggregationCount:
		// суммируем константу 1
		aggExpr = bson.D{{Key: "$sum", Value: 1}}
	case entity.AggregationSum:
		aggExpr = bson.D{{Key: "$sum", Value: "$y"}}
	case entity.AggregationAvg:
		aggExpr = bson.D{{Key: "$avg", Value: "$y"}}
	case entity.AggregationMin:
		aggExpr = bson.D{{Key: "$min", Value: "$y"}}
	case entity.AggregationMax:
		aggExpr = bson.D{{Key: "$max", Value: "$y"}}
	default:
		return nil, "", fmt.Errorf("unsupported aggregation %v", aggType)
	}
	group := bson.D{{"$group", bson.D{{Key: "_id", Value: "$x"}, {Key: "value", Value: aggExpr}}}}
	pipe = append(pipe, group)

	// $project: превращаем _id в строку и переименовываем в x
	pipe = append(pipe, bson.D{{Key: "$project", Value: bson.D{
		{Key: "x", Value: bson.D{{Key: "$toString", Value: "$_id"}}},
		{Key: "value", Value: 1},
	}}})

	return pipe, collName, nil
}

func (r *repository) Aggregate(ctx context.Context, x, y string, aggType entity.Aggregation) ([]entity.Point, error) {
	pipe, collName, err := buildPipeline(x, y, aggType)
	if err != nil {
		return nil, err
	}

	var coll mongoifc.Collection
	switch collName {
	case "orders":
		coll = r.ordersColl
	case "users":
		coll = r.usersColl
	default:
		return nil, fmt.Errorf("invalid collection name: %s", coll)
	}

	cur, err := coll.Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}

	var out []entity.Point
	if err = cur.All(ctx, &out); err != nil {
		return nil, err
	}

	return out, nil
}

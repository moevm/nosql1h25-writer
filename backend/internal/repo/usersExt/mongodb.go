package usersExt

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

const usersCollection = "users"

func (r *mongodbRepo) FindUsers(ctx context.Context, params entity.UserSearchParams) ([]entity.UserExt, int64, error) {
	collection := r.db.Collection(usersCollection)

	filter := bson.M{"active": true}

	if len(params.ProfileFilter) > 0 {

		andConditions := []bson.M{}
		for _, role := range params.ProfileFilter {
			andConditions = append(andConditions, bson.M{"profiles": bson.M{"$elemMatch": bson.M{"role": role}}})
		}
		if len(andConditions) > 0 {
			filter["$and"] = andConditions
		}
	}

	opts := options.Find()
	opts.SetSkip(params.Offset)
	opts.SetLimit(params.Limit)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []entity.UserExt{}, 0, nil
		}

		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []entity.UserExt
	if err = cursor.All(ctx, &users); err != nil {

		return nil, 0, err
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {

		return nil, 0, err
	}

	return users, total, nil
}

package users

import (
	"context"
	"errors"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type repository struct {
	usersColl mongoifc.Collection
}

func New(usersColl mongoifc.Collection) Repo {
	return &repository{usersColl: usersColl}
}

func (r *repository) GetByEmail(ctx context.Context, email string) (u entity.User, _ error) {
	if err := r.usersColl.FindOne(ctx, bson.M{"email": email, "active": true}).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return u, ErrUserNotFound
		}

		return u, err
	}

	return u, nil
}

func (r *repository) GetById(ctx context.Context, ID primitive.ObjectID) (u entity.User, _ error) {
	if err := r.usersColl.FindOne(ctx, bson.M{"_id": ID, "active": true}).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return u, ErrUserNotFound
		}

		return u, err
	}

	return u, nil
}

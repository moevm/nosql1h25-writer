package users

import (
	"context"
	"errors"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

func (r *repository) GetByID(ctx context.Context, id primitive.ObjectID) (u entity.User, _ error) {
	if err := r.usersColl.FindOne(ctx, bson.M{"_id": id, "active": true}).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return u, ErrUserNotFound
		}

		return u, err
	}

	return u, nil
}

func (r *repository) Deposit(ctx context.Context, userID primitive.ObjectID, amount int) (int, error) {
	var u entity.User
	err := r.usersColl.FindOneAndUpdate(
		ctx,
		bson.M{"_id": userID, "active": true},
		bson.M{"$inc": bson.M{"balance": amount}},
		options.FindOneAndUpdate().
			SetReturnDocument(options.After).
			SetProjection(bson.M{"_id": 0, "balance": 1}),
	).Decode(&u)

	if err != nil {
		return 0, err
	}

	return u.Balance, nil
}

func (r *repository) Withdraw(ctx context.Context, userID primitive.ObjectID, amount int) (int, error) {
	var u entity.User
	err := r.usersColl.FindOneAndUpdate(
		ctx,
		bson.M{
			"_id":     userID,
			"active":  true,
			"balance": bson.M{"$gte": amount},
		},
		bson.M{"$inc": bson.M{"balance": -amount}},
		options.FindOneAndUpdate().
			SetReturnDocument(options.After).
			SetProjection(bson.M{"_id": 0, "balance": 1}),
	).Decode(&u)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, ErrInsufficientFunds
		}

		return 0, err
	}

	return u.Balance, nil
}

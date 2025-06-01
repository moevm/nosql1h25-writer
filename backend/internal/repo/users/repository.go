package users

import (
	"context"
	"errors"

	"github.com/AlekSi/pointer"
	"github.com/jonboulle/clockwork"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type repository struct {
	usersColl mongoifc.Collection
	clock     clockwork.Clock
}

func New(usersColl mongoifc.Collection, clock clockwork.Clock) Repo {
	return &repository{
		usersColl: usersColl,
		clock:     clock,
	}
}

func (r *repository) Find(ctx context.Context, in FindIn) (FindOut, error) {
	filter := bson.M{"active": true}

	freelancerRatingCond := bson.M{}
	if in.MinFreelancerRating != nil {
		freelancerRatingCond["$gte"] = *in.MinFreelancerRating
	}
	if in.MaxFreelancerRating != nil {
		freelancerRatingCond["$lte"] = *in.MaxFreelancerRating
	}
	if len(freelancerRatingCond) > 0 {
		filter["freelancer.rating"] = freelancerRatingCond
	}

	clientRatingCond := bson.M{}
	if in.MinClientRating != nil {
		clientRatingCond["$gte"] = *in.MinClientRating
	}
	if in.MaxFreelancerRating != nil {
		clientRatingCond["$lte"] = *in.MaxClientRating
	}
	if len(clientRatingCond) > 0 {
		filter["client.rating"] = clientRatingCond
	}

	createdAtCond := bson.M{}
	if in.MinCreatedAt != nil {
		createdAtCond["$gte"] = *in.MinCreatedAt
	}
	if in.MaxCreatedAt != nil {
		createdAtCond["$lte"] = *in.MaxCreatedAt
	}
	if len(createdAtCond) > 0 {
		filter["createdAt"] = createdAtCond
	}

	balanceCond := bson.M{}
	if in.MinBalance != nil {
		balanceCond["$gte"] = *in.MinBalance
	}
	if in.MaxBalance != nil {
		balanceCond["$lte"] = *in.MaxBalance
	}
	if len(balanceCond) > 0 {
		filter["balance"] = balanceCond
	}

	if pointer.Get(in.NameSearch) != "" {
		filter["displayName"] = bson.M{"$regex": *in.NameSearch, "$options": "i"}
	}
	if pointer.Get(in.EmailSearch) != "" {
		filter["email"] = bson.M{"$regex": *in.EmailSearch, "$options": "i"}
	}

	if len(in.Roles) != 0 {
		filter["systemRole"] = bson.M{"$in": in.Roles}
	}

	var sort bson.D
	switch pointer.Get(in.SortBy) {
	case "newest":
		sort = bson.D{{Key: "createdAt", Value: -1}}
	case "oldest":
		sort = bson.D{{Key: "createdAt", Value: 1}}
	case "rich":
		sort = bson.D{{Key: "balance", Value: -1}}
	case "poor":
		sort = bson.D{{Key: "balance", Value: 1}}
	case "name_asc":
		sort = bson.D{{Key: "displayName", Value: 1}}
	case "name_desc":
		sort = bson.D{{Key: "displayName", Value: -1}}
	case "frelancer_rating_asc":
		sort = bson.D{{Key: "freelancer.rating", Value: 1}}
	case "frelancer_rating_desc":
		sort = bson.D{{Key: "freelancer.rating", Value: -1}}
	case "client_rating_asc":
		sort = bson.D{{Key: "client.rating", Value: 1}}
	case "client_rating_desc":
		sort = bson.D{{Key: "client.rating", Value: -1}}
	default:
		sort = bson.D{{Key: "createdAt", Value: -1}}
	}

	opts := options.Find().SetLimit(int64(in.Limit)).SetSkip(int64(in.Offset)).SetSort(sort)
	cursor, err := r.usersColl.Find(ctx, filter, opts)
	if err != nil {
		return FindOut{}, err
	}

	var users []entity.UserExt
	if err := cursor.All(ctx, &users); err != nil {
		return FindOut{}, err
	}

	total, err := r.usersColl.CountDocuments(ctx, filter)
	if err != nil {
		return FindOut{}, err
	}

	return FindOut{Users: users, Total: int(total)}, nil
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

func (r *repository) GetByIDExt(ctx context.Context, userID primitive.ObjectID) (entity.UserExt, error) {
	var user entity.UserExt
	filter := bson.M{
		"_id":    userID,
		"active": true,
	}

	if err := r.usersColl.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.UserExt{}, ErrUserNotFound
		}

		return entity.UserExt{}, err
	}

	return user, nil
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

func (r *repository) Create(ctx context.Context, in CreateIn) (primitive.ObjectID, error) {
	now := r.clock.Now()
	user := entity.DefaultUser(in.Email, in.Password, in.DisplayName, now, now)

	res, err := r.usersColl.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return primitive.ObjectID{}, ErrUserAlreadyExists
		}

		return primitive.ObjectID{}, err
	}

	return res.InsertedID.(primitive.ObjectID), nil //nolint:forcetypeassert
}

func (r *repository) Update(ctx context.Context, in UpdateIn) error {
	now := r.clock.Now()
	update := bson.M{"updatedAt": now}

	if in.DisplayName != nil {
		update["displayName"] = *in.DisplayName
	}

	if in.FreelancerDescription != nil {
		update["freelancer.description"] = *in.FreelancerDescription
		update["freelancer.updatedAt"] = now
	}

	if in.ClientDescription != nil {
		update["client.description"] = *in.ClientDescription
		update["client.updatedAt"] = now
	}

	err := r.usersColl.FindOneAndUpdate(
		ctx,
		bson.M{"_id": in.UserID, "active": true},
		bson.M{"$set": update},
	).Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

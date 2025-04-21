package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type repository struct {
	sessionsColl mongoifc.Collection
	clock        clockwork.Clock
}

func New(sessionsColl mongoifc.Collection, clock clockwork.Clock) Repo {
	return &repository{
		sessionsColl: sessionsColl,
		clock:        clock,
	}
}

func (r *repository) CreateSession(
	ctx context.Context,
	userID primitive.ObjectID,
	ttl time.Duration,
) (entity.RefreshSession, error) {
	now := r.clock.Now()
	session := entity.RefreshSession{
		RefreshToken: uuid.New(),
		UserID:       userID,
		ExpiresAt:    now.Add(ttl),
		CreatedAt:    now,
	}

	res, err := r.sessionsColl.InsertOne(ctx, session)
	if err != nil {
		return entity.RefreshSession{}, err
	}

	session.ID = res.InsertedID.(primitive.ObjectID) //nolint:forcetypeassert
	return session, nil
}

func (r *repository) GetAndDeleteByToken(ctx context.Context, token uuid.UUID) (s entity.RefreshSession, _ error) {
	if err := r.sessionsColl.FindOneAndDelete(ctx, bson.M{"refreshToken": token}).Decode(&s); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return s, ErrSessionNotFound
		}

		return s, err
	}

	return s, nil
}

func (r *repository) DeleteByToken(ctx context.Context, token uuid.UUID) error {
	res, err := r.sessionsColl.DeleteOne(ctx, bson.M{"refreshToken": token})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrSessionNotFound
	}

	return nil
}

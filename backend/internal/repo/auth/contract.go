package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type Repo interface {
	CreateSession(ctx context.Context, userID primitive.ObjectID, ttl time.Duration) (entity.RefreshSession, error)
	GetAndDeleteByToken(ctx context.Context, token uuid.UUID) (entity.RefreshSession, error)
	DeleteByToken(ctx context.Context, token uuid.UUID) error
}

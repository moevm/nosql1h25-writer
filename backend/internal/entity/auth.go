package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthData struct {
	AccessToken string
	Session     RefreshSession
}

type RefreshSession struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	RefreshToken uuid.UUID          `bson:"refreshToken"`
	UserID       primitive.ObjectID `bson:"userId"`
	ExpiresAt    time.Time          `bson:"expiresAt"`
	CreatedAt    time.Time          `bson:"createdAt"`
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	UserID     primitive.ObjectID `json:"userId"`
	SystemRole SystemRoleType     `json:"systemRole"`
}

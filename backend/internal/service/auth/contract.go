package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type Service interface {
	Login(ctx context.Context, email, password string) (entity.AuthData, error)
	Refresh(ctx context.Context, refreshToken uuid.UUID) (entity.AuthData, error)
	ParseToken(tokenString string) (*entity.AccessTokenClaims, error)
}

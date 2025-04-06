package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	Login(ctx context.Context, email, password string) (entity.AuthData, error)
	Refresh(ctx context.Context, refreshToken uuid.UUID) (entity.AuthData, error)
	Logout(ctx context.Context, refreshToken uuid.UUID) error
	ParseToken(tokenString string) (*entity.AccessTokenClaims, error)
}

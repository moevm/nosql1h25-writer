package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

func (s *service) generateToken(userID primitive.ObjectID, systemRole entity.SystemRoleType) (string, error) {
	now := s.clock.Now()

	return jwt.NewWithClaims(jwt.SigningMethodHS256, &entity.AccessTokenClaims{
		UserID:     userID,
		SystemRole: systemRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}).SignedString([]byte(s.secretKey))
}

func (s *service) ParseToken(tokenString string) (*entity.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.AccessTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(s.secretKey), nil
	}, jwt.WithExpirationRequired(), jwt.WithLeeway(time.Second*3))
	if err != nil {
		log.Errorf("auth.service.ParseToken - jwt.ParseWithClaims: %v", err)
		return nil, ErrCannotAcceptToken
	}

	claims, ok := token.Claims.(*entity.AccessTokenClaims)
	if !ok {
		log.Error("auth.service.ParseToken: unsuccessful cast to custom claims")
		return nil, ErrCannotAcceptToken
	}

	return claims, nil
}

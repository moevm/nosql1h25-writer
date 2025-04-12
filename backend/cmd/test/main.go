package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

func main() {
	issuedAt := lo.Must(time.Parse(time.RFC3339, "2025-04-06T15:00:00Z"))
	userID := lo.Must(primitive.ObjectIDFromHex("507f1f77bcf86cd799439011"))
	systemRole := entity.SystemRoleTypeUser
	secretKey := "secret"

	fmt.Println(lo.Must(jwt.NewWithClaims(jwt.SigningMethodHS256, &entity.AccessTokenClaims{
		UserID:     userID,
		SystemRole: systemRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(time.Minute)),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	}).SignedString([]byte(secretKey))))
}

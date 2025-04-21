package mw

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/auth"
)

var (
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrNotEnoughRights   = errors.New("not enough rights")
)

const (
	UserIDKey     = "userID"
	SystemRoleKey = "systemRole"
)

type AuthMW struct {
	authService auth.Service
}

func NewAuthMW(authService auth.Service) *AuthMW {
	return &AuthMW{authService: authService}
}

// UserIdentity - middleware to check authorization
//
// - Parse `Authorization` header (expected format `Bearer <JWT token>`)
//
// - Sets {"userID": <primitive.ObjectID>, "systemRole": <entity.SystemRoleType>} in `c echo.Context`
func (m *AuthMW) UserIdentity() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := bearerToken(c.Request())
			if err != nil {
				log.Errorf("AuthMW.UserIdentity - bearerToken: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			}

			claims, err := m.authService.ParseToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			c.Set(UserIDKey, claims.UserID)
			c.Set(SystemRoleKey, claims.SystemRole)

			return next(c)
		}
	}
}

// AdminRole - middleware to check user is admin
//
// Expect that value from `c echo.Context` by "systemRole" key is `entity.SystemRoleTypeAdmin`
func (m *AuthMW) Role(targetRole entity.SystemRoleType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get(SystemRoleKey).(entity.SystemRoleType)
			if !ok || role != targetRole {
				return echo.NewHTTPError(http.StatusForbidden, ErrNotEnoughRights.Error())
			}

			return next(c)
		}
	}
}

func bearerToken(req *http.Request) (string, error) {
	const prefix = "Bearer "

	header := req.Header.Get(echo.HeaderAuthorization)

	if len(header) == 0 {
		return "", ErrInvalidAuthHeader
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], nil
	}

	return "", ErrInvalidAuthHeader
}

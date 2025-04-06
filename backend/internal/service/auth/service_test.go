package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	auth_repo_mocks "github.com/moevm/nosql1h25-writer/backend/internal/repo/auth/mocks"
	users_repo_mocks "github.com/moevm/nosql1h25-writer/backend/internal/repo/users/mocks"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/auth"
	"github.com/moevm/nosql1h25-writer/backend/pkg/hasher"
)

func TestService_Login(t *testing.T) {
	var (
		startTime       = lo.Must(time.Parse(time.RFC3339, "2025-04-06T15:00:00Z"))
		ctx             = context.TODO()
		email           = "test@email.ru"
		password        = "qwerty12345"
		secretKey       = "secret"
		accessTokenTTL  = time.Minute
		refreshTokenTTL = time.Hour
	)

	user := entity.User{
		ID:         primitive.NewObjectID(),
		SystemRole: entity.SystemRoleTypeUser,
		Password:   "password",
	}

	token := lo.Must(jwt.NewWithClaims(jwt.SigningMethodHS256, &entity.AccessTokenClaims{
		UserID:     user.ID,
		SystemRole: user.SystemRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(startTime.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(startTime),
		},
	}).SignedString([]byte(secretKey)))

	session := entity.RefreshSession{RefreshToken: uuid.New()}

	type MockBehavior func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo, h *hasher.MockPasswordHasher)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		want         entity.AuthData
		wantErr      error
	}{
		{
			name: "successful",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo, h *hasher.MockPasswordHasher) {
				u.EXPECT().GetByEmail(ctx, email).Return(user, nil)
				h.EXPECT().Match(password, user.Password).Return(true)
				a.EXPECT().CreateSession(ctx, user.ID, refreshTokenTTL).Return(session, nil)
			},
			want: entity.AuthData{
				AccessToken: token,
				Session:     session,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockUsersRepo := users_repo_mocks.NewMockRepo(ctrl)
			mockAuthRepo := auth_repo_mocks.NewMockRepo(ctrl)
			mockPasswordHasher := hasher.NewMockPasswordHasher(ctrl)
			mockClock := clockwork.NewFakeClockAt(startTime)

			tc.mockBehavior(mockUsersRepo, mockAuthRepo, mockPasswordHasher)

			s := auth.New(
				mockUsersRepo,
				mockAuthRepo,
				mockPasswordHasher,
				mockClock,
				secretKey,
				accessTokenTTL,
				refreshTokenTTL,
			)

			got, err := s.Login(ctx, email, password)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

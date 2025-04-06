package auth_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	auth_repo "github.com/moevm/nosql1h25-writer/backend/internal/repo/auth"
	auth_repo_mocks "github.com/moevm/nosql1h25-writer/backend/internal/repo/auth/mocks"
	users_repo "github.com/moevm/nosql1h25-writer/backend/internal/repo/users"
	users_repo_mocks "github.com/moevm/nosql1h25-writer/backend/internal/repo/users/mocks"
	auth_service "github.com/moevm/nosql1h25-writer/backend/internal/service/auth"
	users_service "github.com/moevm/nosql1h25-writer/backend/internal/service/users"
	"github.com/moevm/nosql1h25-writer/backend/pkg/hasher"
)

func TestService_Login(t *testing.T) {
	logrus.SetOutput(io.Discard)
	var (
		startTime       = lo.Must(time.Parse(time.RFC3339, "2025-04-06T15:00:00Z"))
		arbitraryErr    = errors.New("arbitrary error")
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

	for _, tc := range []struct {
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
		{
			name: "user not found",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo, h *hasher.MockPasswordHasher) {
				u.EXPECT().GetByEmail(ctx, email).Return(entity.User{}, users_repo.ErrUserNotFound)
			},
			wantErr: users_service.ErrUserNotFound,
		},
		{
			name: "cannot get user",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo, h *hasher.MockPasswordHasher) {
				u.EXPECT().GetByEmail(ctx, email).Return(entity.User{}, arbitraryErr)
			},
			wantErr: users_service.ErrCannotGetUser,
		},
		{
			name: "wrong password",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo, h *hasher.MockPasswordHasher) {
				u.EXPECT().GetByEmail(ctx, email).Return(user, nil)
				h.EXPECT().Match(password, user.Password).Return(false)
			},
			wantErr: auth_service.ErrWrongPassword,
		},
		{
			name: "cannot create session",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo, h *hasher.MockPasswordHasher) {
				u.EXPECT().GetByEmail(ctx, email).Return(user, nil)
				h.EXPECT().Match(password, user.Password).Return(true)
				a.EXPECT().CreateSession(ctx, user.ID, refreshTokenTTL).Return(entity.RefreshSession{}, arbitraryErr)
			},
			wantErr: auth_service.ErrCannotCreateSession,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockUsersRepo := users_repo_mocks.NewMockRepo(ctrl)
			mockAuthRepo := auth_repo_mocks.NewMockRepo(ctrl)
			mockPasswordHasher := hasher.NewMockPasswordHasher(ctrl)
			mockClock := clockwork.NewFakeClockAt(startTime)

			tc.mockBehavior(mockUsersRepo, mockAuthRepo, mockPasswordHasher)

			s := auth_service.New(
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

func TestService_Refresh(t *testing.T) {
	logrus.SetOutput(io.Discard)
	var (
		startTime       = lo.Must(time.Parse(time.RFC3339, "2025-04-06T15:00:00Z"))
		arbitraryErr    = errors.New("arbitrary error")
		ctx             = context.TODO()
		refreshToken    = uuid.New()
		secretKey       = "secret"
		accessTokenTTL  = time.Minute
		refreshTokenTTL = time.Hour
	)

	user := entity.User{
		ID:         primitive.NewObjectID(),
		SystemRole: entity.SystemRoleTypeUser,
	}

	accessToken := lo.Must(jwt.NewWithClaims(jwt.SigningMethodHS256, &entity.AccessTokenClaims{
		UserID:     user.ID,
		SystemRole: user.SystemRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(startTime.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(startTime),
		},
	}).SignedString([]byte(secretKey)))

	newSession := entity.RefreshSession{RefreshToken: uuid.New()}

	type MockBehavior func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         entity.AuthData
		wantErr      error
	}{
		{
			name: "successful",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo) {
				a.EXPECT().GetAndDeleteByToken(ctx, refreshToken).
					Return(entity.RefreshSession{ExpiresAt: startTime.Add(time.Hour), UserID: user.ID}, nil)
				u.EXPECT().GetByID(ctx, user.ID).Return(user, nil)
				a.EXPECT().CreateSession(ctx, user.ID, refreshTokenTTL).Return(newSession, nil)
			},
			want: entity.AuthData{
				AccessToken: accessToken,
				Session:     newSession,
			},
		},
		{
			name: "session not found",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo) {
				a.EXPECT().GetAndDeleteByToken(ctx, refreshToken).
					Return(entity.RefreshSession{}, auth_repo.ErrSessionNotFound)
			},
			wantErr: auth_service.ErrSessionNotFound,
		},
		{
			name: "cannot get session",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo) {
				a.EXPECT().GetAndDeleteByToken(ctx, refreshToken).
					Return(entity.RefreshSession{}, arbitraryErr)
			},
			wantErr: auth_service.ErrCannotGetSession,
		},
		{
			name: "session expired",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo) {
				a.EXPECT().GetAndDeleteByToken(ctx, refreshToken).
					Return(entity.RefreshSession{ExpiresAt: startTime.Add(-time.Hour), UserID: user.ID}, nil)
			},
			wantErr: auth_service.ErrSessionExpired,
		},
		{
			name: "user not found",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo) {
				a.EXPECT().GetAndDeleteByToken(ctx, refreshToken).
					Return(entity.RefreshSession{ExpiresAt: startTime.Add(time.Hour), UserID: user.ID}, nil)
				u.EXPECT().GetByID(ctx, user.ID).Return(entity.User{}, users_repo.ErrUserNotFound)
			},
			wantErr: users_service.ErrUserNotFound,
		},
		{
			name: "cannot get user",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo) {
				a.EXPECT().GetAndDeleteByToken(ctx, refreshToken).
					Return(entity.RefreshSession{ExpiresAt: startTime.Add(time.Hour), UserID: user.ID}, nil)
				u.EXPECT().GetByID(ctx, user.ID).Return(entity.User{}, arbitraryErr)
			},
			wantErr: users_service.ErrCannotGetUser,
		},
		{
			name: "cannot create session",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo) {
				a.EXPECT().GetAndDeleteByToken(ctx, refreshToken).
					Return(entity.RefreshSession{ExpiresAt: startTime.Add(time.Hour), UserID: user.ID}, nil)
				u.EXPECT().GetByID(ctx, user.ID).Return(user, nil)
				a.EXPECT().CreateSession(ctx, user.ID, refreshTokenTTL).Return(entity.RefreshSession{}, arbitraryErr)
			},
			wantErr: auth_service.ErrCannotCreateSession,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockUsersRepo := users_repo_mocks.NewMockRepo(ctrl)
			mockAuthRepo := auth_repo_mocks.NewMockRepo(ctrl)
			mockPasswordHasher := hasher.NewMockPasswordHasher(ctrl)
			mockClock := clockwork.NewFakeClockAt(startTime)

			tc.mockBehavior(mockUsersRepo, mockAuthRepo)

			s := auth_service.New(
				mockUsersRepo,
				mockAuthRepo,
				mockPasswordHasher,
				mockClock,
				secretKey,
				accessTokenTTL,
				refreshTokenTTL,
			)

			got, err := s.Refresh(ctx, refreshToken)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

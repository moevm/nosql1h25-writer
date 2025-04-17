package auth_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
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

func TestService_Register(t *testing.T) {
	log.SetOutput(io.Discard)
	var (
		ctx             = context.TODO()
		startTime       = lo.Must(time.Parse(time.RFC3339, "2025-04-06T15:00:00Z"))
		secretKey       = "secret"
		accessTokenTTL  = time.Minute
		refreshTokenTTL = time.Hour
		email           = "test@email.ru"
		password        = "qwerty12345"
		userID          = primitive.NewObjectID()
		hashedPassword  = "bombombini gusini"
		serviceIn       = auth_service.RegisterIn{
			DisplayName: "larili larila",
			Email:       email,
			Password:    password,
		}
		repoCreateIn = users_repo.CreateIn{
			DisplayName: serviceIn.DisplayName,
			Email:       email,
			Password:    hashedPassword,
		}
		user = entity.User{
			ID:         userID,
			SystemRole: entity.SystemRoleTypeUser,
			Password:   hashedPassword,
		}
		session = entity.RefreshSession{RefreshToken: uuid.New()}
		token   = lo.Must(jwt.NewWithClaims(jwt.SigningMethodHS256, &entity.AccessTokenClaims{
			UserID:     userID,
			SystemRole: user.SystemRole,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(startTime.Add(accessTokenTTL)),
				IssuedAt:  jwt.NewNumericDate(startTime),
			},
		}).SignedString([]byte(secretKey)))
	)

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
				h.EXPECT().Hash(password).Return(hashedPassword, nil)
				u.EXPECT().Create(ctx, repoCreateIn).Return(userID, nil)
				u.EXPECT().GetByEmail(ctx, email).Return(user, nil)
				h.EXPECT().Match(password, hashedPassword).Return(true)
				a.EXPECT().CreateSession(ctx, userID, refreshTokenTTL).Return(session, nil)
			},
			want: entity.AuthData{
				AccessToken: token,
				Session:     session,
			},
		},
		{
			name: "user already exists",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo, h *hasher.MockPasswordHasher) {
				h.EXPECT().Hash(password).Return(hashedPassword, nil)
				u.EXPECT().Create(ctx, repoCreateIn).Return(primitive.ObjectID{}, users_repo.ErrUserAlreadyExists)
			},
			wantErr: users_service.ErrUserAlreadyExists,
		},
		{
			name: "cannot create user",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo, h *hasher.MockPasswordHasher) {
				h.EXPECT().Hash(password).Return(hashedPassword, nil)
				u.EXPECT().Create(ctx, repoCreateIn).Return(primitive.ObjectID{}, assert.AnError)
			},
			wantErr: users_service.ErrCannotCreateUser,
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

			got, err := s.Register(ctx, serviceIn)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestService_Login(t *testing.T) {
	log.SetOutput(io.Discard)
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
				u.EXPECT().GetByEmail(ctx, email).Return(entity.User{}, assert.AnError)
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
				a.EXPECT().CreateSession(ctx, user.ID, refreshTokenTTL).Return(entity.RefreshSession{}, assert.AnError)
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
	log.SetOutput(io.Discard)
	var (
		startTime       = lo.Must(time.Parse(time.RFC3339, "2025-04-06T15:00:00Z"))
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
					Return(entity.RefreshSession{}, assert.AnError)
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
				u.EXPECT().GetByID(ctx, user.ID).Return(entity.User{}, assert.AnError)
			},
			wantErr: users_service.ErrCannotGetUser,
		},
		{
			name: "cannot create session",
			mockBehavior: func(u *users_repo_mocks.MockRepo, a *auth_repo_mocks.MockRepo) {
				a.EXPECT().GetAndDeleteByToken(ctx, refreshToken).
					Return(entity.RefreshSession{ExpiresAt: startTime.Add(time.Hour), UserID: user.ID}, nil)
				u.EXPECT().GetByID(ctx, user.ID).Return(user, nil)
				a.EXPECT().CreateSession(ctx, user.ID, refreshTokenTTL).Return(entity.RefreshSession{}, assert.AnError)
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

func TestService_Logout(t *testing.T) {
	log.SetOutput(io.Discard)
	var (
		startTime       = lo.Must(time.Parse(time.RFC3339, "2025-04-06T15:00:00Z"))
		ctx             = context.TODO()
		refreshToken    = uuid.New()
		secretKey       = "secret"
		accessTokenTTL  = time.Minute
		refreshTokenTTL = time.Hour
	)

	type MockBehavior func(a *auth_repo_mocks.MockRepo)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "successful",
			mockBehavior: func(a *auth_repo_mocks.MockRepo) {
				a.EXPECT().DeleteByToken(ctx, refreshToken).Return(nil)
			},
		},
		{
			name: "session not found",
			mockBehavior: func(a *auth_repo_mocks.MockRepo) {
				a.EXPECT().DeleteByToken(ctx, refreshToken).Return(auth_repo.ErrSessionNotFound)
			},
			wantErr: auth_service.ErrSessionNotFound,
		},
		{
			name: "cannot delete session",
			mockBehavior: func(a *auth_repo_mocks.MockRepo) {
				a.EXPECT().DeleteByToken(ctx, refreshToken).Return(assert.AnError)
			},
			wantErr: auth_service.ErrCannotDeleteSession,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockUsersRepo := users_repo_mocks.NewMockRepo(ctrl)
			mockAuthRepo := auth_repo_mocks.NewMockRepo(ctrl)
			mockPasswordHasher := hasher.NewMockPasswordHasher(ctrl)
			mockClock := clockwork.NewFakeClockAt(startTime)

			tc.mockBehavior(mockAuthRepo)

			s := auth_service.New(
				mockUsersRepo,
				mockAuthRepo,
				mockPasswordHasher,
				mockClock,
				secretKey,
				accessTokenTTL,
				refreshTokenTTL,
			)

			err := s.Logout(ctx, refreshToken)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestService_ParseToken(t *testing.T) {
	log.SetOutput(io.Discard)
	var (
		issuedTimeValid   = time.Now()
		issuedTimeExpired = issuedTimeValid.Add(-time.Hour)
		userID            = lo.Must(primitive.ObjectIDFromHex("507f1f77bcf86cd799439011"))
		systemRole        = entity.SystemRoleTypeUser
		secretKey         = "secret"
		accessTokenTTL    = time.Minute * 10
		refreshTokenTTL   = time.Hour
	)

	validClaims := entity.AccessTokenClaims{
		UserID:     userID,
		SystemRole: systemRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(issuedTimeValid.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(issuedTimeValid),
		},
	}
	expiredClaims := validClaims
	expiredClaims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(issuedTimeExpired.Add(accessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(issuedTimeExpired),
	}

	validTokenString := lo.Must(jwt.NewWithClaims(jwt.SigningMethodHS256, &validClaims).
		SignedString([]byte(secretKey)))
	diffSecretKeyTokenString := lo.Must(jwt.NewWithClaims(jwt.SigningMethodHS256, &validClaims).
		SignedString([]byte(secretKey + "a")))
	expiredTokenString := lo.Must(jwt.NewWithClaims(jwt.SigningMethodHS256, &expiredClaims).
		SignedString([]byte(secretKey)))
	diffSigningMethodTokenString := lo.Must(jwt.NewWithClaims(jwt.SigningMethodES256, &validClaims).
		SignedString(lo.Must(ecdsa.GenerateKey(elliptic.P256(), rand.Reader))))

	for _, tc := range []struct {
		name              string
		accessTokenString string
		want              *entity.AccessTokenClaims
		wantErr           error
	}{
		{
			name:              "successful",
			accessTokenString: validTokenString,
			want:              &validClaims,
		},
		{
			name:              "diff secret key",
			accessTokenString: diffSecretKeyTokenString,
			wantErr:           auth_service.ErrCannotAcceptToken,
		},
		{
			name:              "expired token",
			accessTokenString: expiredTokenString,
			wantErr:           auth_service.ErrCannotAcceptToken,
		},
		{
			name:              "diff signing method",
			accessTokenString: diffSigningMethodTokenString,
			wantErr:           auth_service.ErrCannotAcceptToken,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockUsersRepo := users_repo_mocks.NewMockRepo(ctrl)
			mockAuthRepo := auth_repo_mocks.NewMockRepo(ctrl)
			mockPasswordHasher := hasher.NewMockPasswordHasher(ctrl)
			mockClock := clockwork.NewFakeClock()

			s := auth_service.New(
				mockUsersRepo,
				mockAuthRepo,
				mockPasswordHasher,
				mockClock,
				secretKey,
				accessTokenTTL,
				refreshTokenTTL,
			)

			got, err := s.ParseToken(tc.accessTokenString)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

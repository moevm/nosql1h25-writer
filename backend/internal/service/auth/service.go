package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	log "github.com/sirupsen/logrus"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/repo/auth"
	users_repo "github.com/moevm/nosql1h25-writer/backend/internal/repo/users"
	users_service "github.com/moevm/nosql1h25-writer/backend/internal/service/users"
	"github.com/moevm/nosql1h25-writer/backend/pkg/hasher"
)

type service struct {
	usersRepo       users_repo.Repo
	authRepo        auth.Repo
	hasher          hasher.PasswordHasher
	clock           clockwork.Clock
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func New(
	usersRepo users_repo.Repo,
	authRepo auth.Repo,
	hasher hasher.PasswordHasher,
	clock clockwork.Clock,
	secretKey string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) Service {
	return &service{
		usersRepo:       usersRepo,
		authRepo:        authRepo,
		hasher:          hasher,
		clock:           clock,
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (entity.AuthData, error) {
	user, err := s.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, users_repo.ErrUserNotFound) {
			return entity.AuthData{}, users_service.ErrUserNotFound
		}

		log.Errorf("auth.service.Login - s.usersRepo.GetByEmail: %v", err)
		return entity.AuthData{}, users_service.ErrCannotGetUser
	}

	if !s.hasher.Match(password, user.Password) {
		return entity.AuthData{}, ErrWrongPassword
	}

	accessToken, err := s.generateToken(user.ID, user.SystemRole)
	if err != nil {
		log.Errorf("auth.service.Login - s.generateToken: %v", err)
		return entity.AuthData{}, ErrCannotGenerateToken
	}

	session, err := s.authRepo.CreateSession(ctx, user.ID, s.refreshTokenTTL)
	if err != nil {
		log.Errorf("auth.service.Login - s.authRepo.CreateSession: %v", err)
		return entity.AuthData{}, ErrCannotCreateSession
	}

	return entity.AuthData{
		AccessToken: accessToken,
		Session:     session,
	}, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken uuid.UUID) (entity.AuthData, error) {
	curSession, err := s.authRepo.GetAndDeleteByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrSessionNotFound) {
			return entity.AuthData{}, ErrSessionNotFound
		}

		log.Errorf("auth.service.Refresh - s.authRepo.GetAndDeleteByToken: %v", err)
		return entity.AuthData{}, ErrCannotGetSession
	}

	if s.clock.Now().After(curSession.ExpiresAt) {
		return entity.AuthData{}, ErrSessionExpired
	}

	user, err := s.usersRepo.GetByID(ctx, curSession.UserID)
	if err != nil {
		if errors.Is(err, users_repo.ErrUserNotFound) {
			return entity.AuthData{}, users_service.ErrUserNotFound
		}

		log.Errorf("auth.service.Refresh - s.usersRepo.GetById: %v", err)
		return entity.AuthData{}, users_service.ErrCannotGetUser
	}

	accessToken, err := s.generateToken(user.ID, user.SystemRole)
	if err != nil {
		log.Errorf("auth.service.Refresh - s.generateToken: %v", err)
		return entity.AuthData{}, ErrCannotGenerateToken
	}

	newSession, err := s.authRepo.CreateSession(ctx, user.ID, s.refreshTokenTTL)
	if err != nil {
		log.Errorf("auth.service.Refresh - s.authRepo.CreateSession: %v", err)
		return entity.AuthData{}, ErrCannotCreateSession
	}

	return entity.AuthData{
		AccessToken: accessToken,
		Session:     newSession,
	}, nil
}

func (s *service) Logout(ctx context.Context, refreshToken uuid.UUID) error {
	if err := s.authRepo.DeleteByToken(ctx, refreshToken); err != nil {
		if errors.Is(err, auth.ErrSessionNotFound) {
			return ErrSessionNotFound
		}

		log.Errorf("auth.service.Logout - s.authRepo.DeleteByToken: %v", err)
		return ErrCannotDeleteSession
	}

	return nil
}

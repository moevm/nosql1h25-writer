package usersExt

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/repo/usersExt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo usersExt.Repo
}

func New(repo usersExt.Repo) Service {
	return &service{repo: repo}
}

func (s *service) FindUsers(ctx context.Context, params entity.UserSearchParams) ([]entity.UserExt, int64, error) {
	log.WithFields(log.Fields{
		"offset":  params.Offset,
		"limit":   params.Limit,
		"profile": params.ProfileFilter,
	}).Debug("Service: Calling repo.FindUsers")

	users, total, err := s.repo.FindUsers(ctx, params)
	if err != nil {
		log.WithError(err).Error("Service: Error finding users in repo")
		return nil, 0, err
	}

	log.WithFields(log.Fields{
		"found_count": len(users),
		"total_count": total,
	}).Debug("Service: Found users in repo")

	return users, total, nil
}

func (s *service) FindUserByID(ctx context.Context, userID primitive.ObjectID) (entity.UserExt, error) {
	log.WithField("userID", userID.Hex()).Debug("Service: Calling repo.FindUserByID")

	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		if !errors.Is(err, usersExt.ErrUserNotFound) {
			log.WithError(err).WithField("userID", userID.Hex()).Error("Service: Error finding user by ID in repo")
		}
		return entity.UserExt{}, err
	}

	log.WithField("userID", userID.Hex()).Debug("Service: Found user by ID in repo")
	return user, nil
}

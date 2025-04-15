package usersExt

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/repo/usersExt"
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

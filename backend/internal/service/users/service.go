package users

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type service struct {
}

func New( /* repo repository.Repository */ ) Service {
	return &service{}
}

func (s *service) FindUsers(ctx context.Context, params entity.UserSearchParams) ([]entity.UserExt, int64, error) {
	log.WithFields(log.Fields{
		"offset":  params.Offset,
		"limit":   params.Limit,
		"profile": params.ProfileFilter,
	}).Info("Service: FindUsers called (stub)")

	return []entity.UserExt{}, 0, nil
}

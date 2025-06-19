package stats

import (
	"context"

	"github.com/moevm/nosql1h25-writer/backend/internal/repo/stats"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type service struct {
	statsRepo    stats.Repo
	usersService users.Service
}

func New(statsRepo stats.Repo, usersService users.Service) Service {
	return &service{
		statsRepo:    statsRepo,
		usersService: usersService,
	}
}

func (s *service) Graph(ctx context.Context, in GraphIn) ([]GraphOut, error) {
	return []GraphOut{}, nil
}

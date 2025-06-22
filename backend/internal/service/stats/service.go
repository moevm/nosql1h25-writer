package stats

import (
	"context"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/repo/stats"
)

type service struct {
	statsRepo stats.Repo
}

func New(statsRepo stats.Repo) Service {
	return &service{
		statsRepo: statsRepo,
	}
}

func (s *service) Graph(ctx context.Context, in GraphIn) ([]GraphOut, error) {
	if !validRequest(in.X, in.Y, in.AggType) {
		return nil, ErrInvalidRequest
	}

	points, err := s.statsRepo.Aggregate(ctx, in.X, in.Y, in.AggType)
	if err != nil {
		log.Errorf("s.Graph - s.statsRepo.Aggregate: %v", err)
		return nil, ErrCannotAggregate
	}

	return lo.Map(points, func(item entity.Point, _ int) GraphOut {
		return GraphOut{
			X: item.X,
			Y: item.Y,
		}
	}), nil
}

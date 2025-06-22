package stats

import (
	"context"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Repo
type Repo interface {
	Aggregate(ctx context.Context, x, y string, aggType entity.Aggregation) ([]entity.Point, error)
}

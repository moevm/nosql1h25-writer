package stats

import (
	"context"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	Graph(ctx context.Context, in GraphIn) ([]GraphOut, error)
}

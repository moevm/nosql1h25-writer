package orders

import (
	"context"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	Find(ctx context.Context, offset, limit int) (FindOut, error)
	// FindOrders(ctx context.Context, offset, limit int64) ([]orders.FindOrdersOut, int64, error)
	// GetByID(ctx context.Context, id string) (orders.FindOrdersOut, error)
}

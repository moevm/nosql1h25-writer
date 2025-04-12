package orders

import (
	"context"

	"github.com/moevm/nosql1h25-writer/backend/internal/repo/orders"
)

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	FindOrders(ctx context.Context, offset, limit int64) ([]orders.FindOrdersOut, int64, error)
	GetByID(ctx context.Context, id string) (orders.FindOrdersOut, error)
}

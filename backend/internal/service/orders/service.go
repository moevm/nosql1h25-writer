package orders

import "github.com/moevm/nosql1h25-writer/backend/internal/repo/orders"

type service struct {
	ordersRepo orders.Repo
}

func New(ordersRepo orders.Repo) Service {
	return &service{ordersRepo: ordersRepo}
}

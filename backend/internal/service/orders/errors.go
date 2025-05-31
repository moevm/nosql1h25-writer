package orders

import "errors"

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrCannotGetOrder    = errors.New("cannot get order")
	ErrCannotCreateOrder = errors.New("cannot create order")
	ErrCannotUpdateOrder = errors.New("cannot update order")
	ErrCannotFindOrders  = errors.New("cannot find orders")
	ErrCannotResponse    = errors.New("cannot response order")
)

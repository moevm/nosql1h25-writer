package orders

import "errors"

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrCannotGetOrder    = errors.New("cannot get order")
	ErrCannotFindOrders  = errors.New("cannot find orders")
	ErrCannotCreateOrder = errors.New("cannot create order")
)

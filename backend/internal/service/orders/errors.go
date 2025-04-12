package orders

import "errors"

var (
	ErrOrdersNotFound    = errors.New("orders not found")
	ErrCannotGetOrders   = errors.New("cannot get orders")
	ErrInvalidPagination = errors.New("invalid pagination parameters")
)

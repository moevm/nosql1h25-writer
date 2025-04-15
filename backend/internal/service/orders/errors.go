package orders

import "errors"

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrCannotGetOrder = errors.New("cannot get order")
)

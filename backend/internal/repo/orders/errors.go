package orders

import "errors"

var ErrNotFound = errors.New("orders not found")
var ErrCannotGetOrders = errors.New("cannot get orders")
var ErrInvalidPagination = errors.New("invalid pagination parameters")

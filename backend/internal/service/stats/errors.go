package stats

import "errors"

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrCannotAggregate = errors.New("cannot aggregate")
)

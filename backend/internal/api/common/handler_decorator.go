package common

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
)

type handler[T any] interface {
	Handle(c echo.Context, in T) error
}

type bindAndValidateDecorator[T any] struct {
	inner handler[T]
}

func NewBindAndValidate[T any](inner handler[T]) api.Handler {
	return &bindAndValidateDecorator[T]{inner: inner}
}

func (d *bindAndValidateDecorator[T]) Handle(c echo.Context) error {
	var in T

	if err := c.Bind(&in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return d.inner.Handle(c, in)
}

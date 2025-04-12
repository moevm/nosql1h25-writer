package users

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrCannotGetUser     = errors.New("cannot get user")
	ErrUpdateBalance     = errors.New("cannot update balance")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

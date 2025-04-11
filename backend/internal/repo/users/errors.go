package users

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrCannotDeposit    = errors.New("cannot deposit")
	ErrCannotWithdraw   = errors.New("cannot withdraw")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
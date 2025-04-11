package users

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrCannotGetUser = errors.New("cannot get user")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrCannotDeposit     = errors.New("cannot deposit")
	ErrCannotWithdraw    = errors.New("cannot withdraw")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

package users

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrCannotGetUser     = errors.New("cannot get user")
	ErrUpdateBalance     = errors.New("cannot update balance")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCannotCreateUser  = errors.New("cannot create user")
)

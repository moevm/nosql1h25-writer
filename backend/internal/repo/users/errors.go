package users

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCannotUpdateUser  = errors.New("cannot update user")
)

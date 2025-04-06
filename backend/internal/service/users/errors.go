package users

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrCannotGetUser = errors.New("cannot get user")
)

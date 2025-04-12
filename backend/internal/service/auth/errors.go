package auth

import "errors"

var (
	ErrWrongPassword = errors.New("wrong password")

	ErrCannotGenerateToken = errors.New("cannot generate token")
	ErrCannotAcceptToken   = errors.New("cannot accept token")

	ErrSessionNotFound     = errors.New("session not found")
	ErrSessionExpired      = errors.New("session expired")
	ErrCannotGetSession    = errors.New("cannot get session")
	ErrCannotCreateSession = errors.New("cannot create session")
	ErrCannotDeleteSession = errors.New("cannot delete session")
)

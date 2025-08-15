package user

import "errors"

var (
	ErrInvalidNickname     = errors.New("invalid nickname")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrExistNickname       = errors.New("user with this nickname exists")
	ErrExistEmail          = errors.New("user with this email exists")
	ErrWrongLengthNickname = errors.New("user nickname length is wrong")
	ErrWrongLengthEmail    = errors.New("user email length is wrong")

	ErrNotFound = errors.New("user not found")
)

package user

import "errors"

var (
	ErrExistNickname       = errors.New("user with this nickname exists")
	ErrExistEmail          = errors.New("user with this email exists")
	ErrWrongLengthNickname = errors.New("user nickname length is wrong")
	ErrWrongLengthEmail    = errors.New("user email length is wrong")

	ErrNotFound = errors.New("user not found")
)

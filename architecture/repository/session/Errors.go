package session

import "errors"

var (
	ErrSessionExists = errors.New("session with this user_id exists")
	ErrNotFound      = errors.New("session not found")
)

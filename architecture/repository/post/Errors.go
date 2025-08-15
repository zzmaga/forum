package post

import "errors"

var (
	// ErrInvalidTitle       = errors.New("invalid title")
	ErrInvalidTitleLength = errors.New("invalid title length")
	// ErrInvalidContentLength = errors.New("invalid content length")

	ErrNotFound = errors.New("post not found")
)

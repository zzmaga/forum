package post

import "errors"

var (
	ErrInvalidTitleLength   = errors.New("invalid title length")
	ErrInvalidContentLength = errors.New("invalid content length")
	ErrNotFound             = errors.New("post not found")
)

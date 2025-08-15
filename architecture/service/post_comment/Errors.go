package post_comment

import "fmt"

var (
	ErrInvalidContentLength = fmt.Errorf("invalid content length")
	ErrNotFound             = fmt.Errorf("post_comment not found")
)

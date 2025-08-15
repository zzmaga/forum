package post_comment_vote

import "fmt"

var (
	ErrExists      = fmt.Errorf("post comment vote exists")
	ErrNotFound    = fmt.Errorf("not found")
	ErrInvalidVote = fmt.Errorf("invalid vote")
)

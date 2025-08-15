package post_vote

import "fmt"

var (
	ErrNotFound    = fmt.Errorf("not found")
	ErrInvalidVote = fmt.Errorf("invalid vote")
)

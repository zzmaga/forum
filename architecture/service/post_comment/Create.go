package post_comment

import (
	"fmt"
	"time"

	"forum/architecture/models"
)

func (c *PostCommentService) Create(comment *models.PostComment) (int64, error) {
	comment.Prepare()

	if comment.ValidateContent() != nil {
		return -1, ErrInvalidContentLength
	}

	comment.CreatedAt = time.Now()
	commentId, err := c.repo.Create(comment)
	switch {
	case err == nil:
	case err != nil:
		return -1, fmt.Errorf("c.repo.Create: %w", err)
	}
	return commentId, nil
}

package post_comment

import (
	"fmt"

	"forum/architecture/models"
)

func (c *PostCommentService) GetAllByPostID(postId, offset, limit int64) ([]*models.PostComment, error) {
	comments, err := c.repo.GetAllByPostID(postId, offset, limit)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("c.repo.GetAllByPostID: %w", err)
	}
	return comments, nil
}

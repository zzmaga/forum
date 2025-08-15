package post_comment

import (
	"errors"
	"fmt"

	"forum/architecture/models"
	"forum/architecture/repository/post_comment"
)

func (c *PostCommentService) GetByID(id int64) (*models.PostComment, error) {
	comment, err := c.repo.GetByID(id)
	switch {
	case err == nil:
	case errors.Is(err, post_comment.ErrNotFound):
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("c.repo.GetByID: %w", err)
	}
	return comment, nil
}

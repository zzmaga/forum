package post_comment

import (
	"fmt"

	"forum/architecture/models"
)

func (c *PostCommentRepo) Create(comment *models.PostComment) (int64, error) {
	strCreatedAt := comment.CreatedAt.Format(models.TimeFormat)
	row := c.db.QueryRow(`
INSERT INTO posts_comments (content, user_id, post_id, created_at) VALUES
(?, ?, ?, ?) RETURNING id`, comment.Content, comment.UserId, comment.PostId, strCreatedAt)

	err := row.Scan(&comment.Id)
	switch {
	case err == nil:
	// case strings.HasPrefix(err.Error(), "FOREIGN KEY constraint failed"):
	// 	return -1, ErrNotFound
	case err != nil:
		return -1, fmt.Errorf("row.Scan: %w", err)
	}
	return comment.Id, nil
}

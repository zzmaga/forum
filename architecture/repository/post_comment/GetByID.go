package post_comment

import (
	"fmt"
	"strings"
	"time"

	"forum/architecture/models"
)

func (c *PostCommentRepo) GetByID(id int64) (*models.PostComment, error) {
	row := c.db.QueryRow(`
	SELECT id, content, user_id, post_id, created_at FROM posts_comments
	WHERE id = ?`, id)

	comment := &models.PostComment{}
	var strCreatedAt string

	err := row.Scan(&comment.Id, &comment.Content, &comment.UserId, &comment.PostId, &strCreatedAt)
	switch {
	case err == nil:
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	timeCreatedAt, err := time.ParseInLocation(models.TimeFormat, strCreatedAt, time.Local)
	if err != nil {
		return nil, fmt.Errorf("time.Parse: %w", err)
	}
	comment.CreatedAt = timeCreatedAt

	return comment, nil
}

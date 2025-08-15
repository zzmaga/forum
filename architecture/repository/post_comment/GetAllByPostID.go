package post_comment

import (
	"fmt"
	"time"

	"forum/architecture/models"
)

func (c *PostCommentRepo) GetAllByPostID(postId, offset, limit int64) ([]*models.PostComment, error) {
	if limit == 0 {
		limit = -1
	}

	rows, err := c.db.Query(`
SELECT id, content, user_id, post_id, created_at FROM posts_comments
WHERE post_id = ?
LIMIT ? OFFSET ? 
	`, postId, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("p.db.Query: %w", err)
	}

	comments := []*models.PostComment{}
	for rows.Next() {
		var strCreatedAt string
		comment := &models.PostComment{}

		err := rows.Scan(&comment.Id, &comment.Content, &comment.UserId, &comment.PostId, &strCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		timeCreatedAt, err := time.ParseInLocation(models.TimeFormat, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		comment.CreatedAt = timeCreatedAt
		comments = append(comments, comment)
	}
	return comments, nil
}

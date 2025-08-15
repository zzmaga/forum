package post

import (
	"fmt"
	"forum/architecture/models"
	"strings"
)

func (p *PostRepo) Create(post *models.Post) (int64, error) {
	strCreatedAt := post.CreatedAt.Format(models.TimeFormat)
	row := p.db.QueryRow(`
INSERT INTO posts (title, content, user_id, created_at, updated_at) VALUES
(?, ?, ?, ?, ?) RETURNING id`, post.Title, post.Content, post.UserId, strCreatedAt, strCreatedAt)

	err := row.Scan(&post.Id)
	switch {
	case err == nil:
		return post.Id, nil
	case strings.HasPrefix(err.Error(), "CHECK constraint failed"):
		// Create Error
		switch {
		case strings.Contains(err.Error(), "title"):
			switch {
			case strings.Contains(err.Error(), "LENGTH"):
				return -1, ErrInvalidTitleLength
			}
		}
	}
	return -1, fmt.Errorf("row.Scan: %w", err)
}

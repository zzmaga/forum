package post

import (
	"fmt"
	"strings"
	"time"

	"forum/architecture/models"
)

func (p *PostRepo) GetByID(id int64) (*models.Post, error) {
	row := p.db.QueryRow(`
	SELECT id, title, content, user_id, created_at, updated_at FROM posts
	WHERE id = ?`, id)

	post := &models.Post{}
	var strCreatedAt, strUpdatedAt string
	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &strCreatedAt, &strUpdatedAt)
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
	post.CreatedAt = timeCreatedAt

	timeUpdatedAt, err := time.ParseInLocation(models.TimeFormat, strUpdatedAt, time.Local)
	if err != nil {
		return nil, fmt.Errorf("time.Parse: %w", err)
	}
	post.UpdatedAt = timeUpdatedAt
	return post, nil
}

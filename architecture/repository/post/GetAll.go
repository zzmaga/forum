package post

import (
	"fmt"
	"time"

	"forum/architecture/models"
)

func (p *PostRepo) GetAll(offset, limit int64) ([]*models.Post, error) {
	if limit == 0 {
		limit = -1
	}

	rows, err := p.db.Query(`
SELECT id, title, content, user_id, created_at, updated_at FROM posts
LIMIT ? OFFSET ? 
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("p.db.Query: %w", err)
	}

	posts := []*models.Post{}
	for rows.Next() {
		var strCreatedAt, strUpdatedAt string
		post := &models.Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &strCreatedAt, &strUpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
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

		posts = append(posts, post)
	}

	return posts, nil
}

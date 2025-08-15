package category

import (
	"fmt"
	"time"

	"forum/architecture/models"
)

func (c *CategoryRepo) GetByPostID(postId int64) ([]*models.Category, error) {
	rows, err := c.db.Query(`
SELECT c.id, c.name, c.created_at FROM posts_categories pc
JOIN categories c ON pc.category_id = c.id
WHERE pc.post_id = ?`, postId)
	if err != nil {
		return nil, fmt.Errorf("c.db.Query: %w", err)
	}

	categories := []*models.Category{}
	for rows.Next() {
		var strCreatedAt string
		category := &models.Category{}
		err = rows.Scan(&category.Id, &category.Name, &strCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		timeCreatedAt, err := time.ParseInLocation(models.TimeFormat, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		category.CreatedAt = timeCreatedAt
		categories = append(categories, category)
	}
	return categories, nil
}

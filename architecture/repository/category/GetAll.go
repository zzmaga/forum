package category

import (
	"fmt"
	"time"

	"forum/architecture/models"
)

func (p *CategoryRepo) GetAll(offset, limit int64) ([]*models.Category, error) {
	if limit == 0 {
		limit = -1
	}

	rows, err := p.db.Query(`
SELECT id, name, created_at FROM categories
LIMIT ? OFFSET ? 
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("p.db.Query: %w", err)
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

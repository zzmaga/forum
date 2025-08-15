package category

import (
	"fmt"
	"strings"

	"forum/architecture/models"
)

func (c *CategoryRepo) Create(category *models.Category) (int64, error) {
	strCreatedAt := category.CreatedAt.Format(models.TimeFormat)
	row := c.db.QueryRow(`
	INSERT INTO categories (name, created_at) VALUES
	(?, ?) RETURNING id`, category.Name, strCreatedAt)
	err := row.Scan(&category.Id)

	switch {
	case err == nil:
		return category.Id, nil
	case strings.HasPrefix(err.Error(), "UNIQUE constraint failed"):
		return -1, ErrExistName
	case strings.HasPrefix(err.Error(), "CHECK constraint failed"):
		return -1, ErrCheckLengthName
	}
	return -1, fmt.Errorf("row.Scan: %w", err)
}

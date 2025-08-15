package category

import (
	"fmt"

	"forum/architecture/models"
)

func (c *CategoryService) GetByPostID(postId int64) ([]*models.Category, error) {
	categories, err := c.repo.GetByPostID(postId)
	switch {
	case err == nil:
		return categories, nil
	}
	return nil, fmt.Errorf("c.repo.GetByPostID: %w", err)
}

package category

import (
	"fmt"

	"forum/architecture/models"
)

func (c *CategoryService) GetByNames(names []string) ([]*models.Category, error) {
	cats, err := c.repo.GetByNames(names)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("GetByNames: %w", err)
	}
	return cats, nil
}

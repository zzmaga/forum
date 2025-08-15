package category

import (
	"fmt"

	"forum/architecture/models"
)

func (p *CategoryService) GetAll(offset, limit int64) ([]*models.Category, error) {
	categories, err := p.repo.GetAll(offset, limit)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("p.repo.GetAll: %w", err)
	}
	return categories, nil
}

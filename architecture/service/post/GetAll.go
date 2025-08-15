package post

import (
	"fmt"

	"forum/architecture/models"
)

func (p *PostService) GetAll(offset, limit int64) ([]*models.Post, error) {
	posts, err := p.repo.GetAll(offset, limit)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("p.repo.GetAll: %w", err)
	}
	return posts, nil
}

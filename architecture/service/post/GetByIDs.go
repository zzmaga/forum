package post

import (
	"fmt"

	"forum/architecture/models"
)

func (p *PostService) GetByIDs(ids []int64) ([]*models.Post, error) {
	posts, err := p.repo.GetByIDs(ids)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("GetByIDs: %w", err)
	}
	return posts, nil
}

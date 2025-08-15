package post

import (
	"fmt"

	"forum/architecture/models"
)

func (p *PostService) GetByUserID(userId, offset, limit int64) ([]*models.Post, error) {
	posts, err := p.repo.GetByUserID(userId, offset, limit)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("p.repo.GetByUserID: %w", err)
	}
	return posts, nil
}

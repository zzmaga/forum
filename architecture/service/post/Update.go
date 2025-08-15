package post

import (
	"fmt"
	"time"

	model "forum/architecture/models"
)

func (p *PostService) Update(post *model.Post) error {
	post.Prepare()

	if post.ValidateTitle() != nil {
		return ErrInvalidTitleLength
	} else if post.ValidateContent() != nil {
		return ErrInvalidContentLength
	}

	post.UpdatedAt = time.Now()
	err := p.repo.Update(post)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("p.repo.Update: %w", err)
	}
	return nil
}

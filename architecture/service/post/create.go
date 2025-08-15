package post

import (
	"errors"
	"fmt"
	"time"

	model "forum/architecture/models"
	rpost "forum/architecture/repository/post"
)

func (p *PostService) Create(post *model.Post) (int64, error) {
	post.Prepare()

	if post.ValidateTitle() != nil {
		return -1, ErrInvalidTitleLength
	} else if post.ValidateContent() != nil {
		return -1, ErrInvalidContentLength
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = post.CreatedAt

	postId, err := p.repo.Create(post)
	switch {
	case err == nil:
	case errors.Is(err, rpost.ErrInvalidTitleLength):
		return -1, ErrInvalidTitleLength
	case err != nil:
		return -1, fmt.Errorf("p.repo.Create: %w", err)
	}
	return postId, nil

}

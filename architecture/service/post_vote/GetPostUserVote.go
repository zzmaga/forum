package post_vote

import (
	"errors"
	"fmt"

	"forum/architecture/models"
	"forum/architecture/repository/post_vote"
)

func (p *PostVoteService) GetPostUserVote(userId, postId int64) (*models.PostVote, error) {
	pVote, err := p.repo.GetPostUserVote(userId, postId)
	switch {
	case err == nil:
	case errors.Is(err, post_vote.ErrNotFound):
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("p.repo.GetPostUserVote: %w", err)
	}
	return pVote, nil
}

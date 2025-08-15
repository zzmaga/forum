package post_comment_vote

import (
	"errors"
	"fmt"
	"time"

	"forum/architecture/models"
	"forum/architecture/repository/post_comment_vote"
)

func (p *PostCommentVoteService) Record(vote *models.PostCommentVote) error {
	if vote.Vote < -1 || 1 < vote.Vote {
		return ErrInvalidVote
	}

	vote.CreatedAt = time.Now()
	_, err := p.repo.Create(vote)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, post_comment_vote.ErrExists):
	case errors.Is(err, post_comment_vote.ErrNotFound):
		return ErrNotFound
	case err != nil:
		return fmt.Errorf("p.repo.Create: %w", err)
	}

	vote.UpdatedAt = time.Now()
	err = p.repo.Update(vote)
	switch {
	case err == nil:
	case errors.Is(err, post_comment_vote.ErrNotFound):
		return ErrNotFound
	case err != nil:
		return fmt.Errorf("p.repo.Update: %w", err)
	}
	return nil
}

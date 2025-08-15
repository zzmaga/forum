package post_comment_vote

import (
	"errors"
	"fmt"

	"forum/architecture/models"
	"forum/architecture/repository/post_comment_vote"
)

func (c *PostCommentVoteService) GetCommentUserVote(userId, commentId int64) (*models.PostCommentVote, error) {
	pVote, err := c.repo.GetCommentUserVote(userId, commentId)
	switch {
	case err == nil:
	case errors.Is(err, post_comment_vote.ErrNotFound):
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("repo.GetCommentUserVote: %w", err)
	}
	return pVote, nil
}

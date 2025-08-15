package post_comment_vote

import (
	"fmt"
	"strings"

	"forum/architecture/models"
)

func (p *PostCommentVoteRepo) Update(vote *models.PostCommentVote) error {
	strUpdatedAt := vote.UpdatedAt.Format(models.TimeFormat)
	row := p.db.QueryRow(`
UPDATE posts_comments_votes
SET vote = ?, updated_at = ?
WHERE user_id = ? AND comment_id = ? 
RETURNING id`, vote.Vote, strUpdatedAt, vote.UserId, vote.CommentId)

	err := row.Scan(&vote.Id)
	switch {
	case err == nil:
	case strings.HasPrefix(err.Error(), "FOREIGN KEY constraint failed"):
		return ErrNotFound
	case err != nil:
		return fmt.Errorf("row.Scan: %w", err)
	}
	return nil
}

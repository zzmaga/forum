package post_vote

import (
	"fmt"
	"strings"

	"forum/architecture/models"
)

func (p *PostVoteRepo) Update(vote *models.PostVote) error {
	strUpdatedAt := vote.UpdatedAt.Format(models.TimeFormat)
	row := p.db.QueryRow(`
UPDATE posts_votes
SET vote = ?, updated_at = ?
WHERE user_id = ? AND post_id = ? 
RETURNING id`, vote.Vote, strUpdatedAt, vote.UserId, vote.PostId)

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

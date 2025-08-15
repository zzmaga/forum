package post_comment_vote

import (
	"fmt"
	"strings"

	"forum/architecture/models"
)

func (p *PostCommentVoteRepo) Create(vote *models.PostCommentVote) (int64, error) {
	strCreatedAt := vote.CreatedAt.Format(models.TimeFormat)
	row := p.db.QueryRow(`
INSERT INTO posts_comments_votes (vote, user_id, comment_id, created_at, updated_at) VALUES
(?, ?, ?, ?, ?) RETURNING id`, vote.Vote, vote.UserId, vote.CommentId, strCreatedAt, strCreatedAt)

	err := row.Scan(&vote.Id)
	switch {
	case err == nil:
	case err != nil:
		switch {
		case strings.HasPrefix(err.Error(), "UNIQUE constraint failed"):
			return -1, ErrExists
		case strings.HasPrefix(err.Error(), "FOREIGN KEY constraint failed"):
			return -1, ErrNotFound
		}
		return -1, fmt.Errorf("row.Scan: %w", err)
	}
	return vote.Id, nil
}

package post_vote

import (
	"fmt"
	"strings"
	"time"

	"forum/architecture/models"
)

func (p *PostVoteRepo) GetPostUserVote(userId, postId int64) (*models.PostVote, error) {
	row := p.db.QueryRow(`
	SELECT id, post_id, user_id, vote, created_at FROM posts_votes
	WHERE post_id = ? AND user_id = ?`, postId, userId)
	postVote := &models.PostVote{}

	strCreatedAt := ""
	err := row.Scan(&postVote.Id, &postVote.PostId, &postVote.UserId, &postVote.Vote, &strCreatedAt)

	switch {
	case err == nil:
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	timeCreatedAt, err := time.ParseInLocation(models.TimeFormat, strCreatedAt, time.Local)
	if err != nil {
		return nil, fmt.Errorf("time.Parse: %w", err)
	}
	postVote.CreatedAt = timeCreatedAt
	return postVote, nil
}

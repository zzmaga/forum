package post_comment_vote

import (
	"fmt"
	"strings"
	"time"

	"forum/architecture/models"
)

func (c *PostCommentVoteRepo) GetCommentUserVote(userId, commentId int64) (*models.PostCommentVote, error) {
	row := c.db.QueryRow(`
	SELECT id, comment_id, user_id, vote, created_at FROM posts_comments_votes
	WHERE comment_id = ? AND user_id = ?`, commentId, userId)
	commentVote := &models.PostCommentVote{}

	strCreatedAt := ""
	err := row.Scan(&commentVote.Id, &commentVote.CommentId, &commentVote.UserId, &commentVote.Vote, &strCreatedAt)

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
	commentVote.CreatedAt = timeCreatedAt
	return commentVote, nil
}

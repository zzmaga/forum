package post_comment_vote

import "fmt"

func (c *PostCommentVoteRepo) GetByCommentID(commentId int64) (int64, int64, error) {
	row := c.db.QueryRow(`
SELECT 
	(SELECT COUNT(vote) FROM posts_comments_votes WHERE comment_id = ? AND vote == 1) as up,
    (SELECT COUNT(vote) FROM posts_comments_votes WHERE comment_id = ? AND vote == -1) as down;
	`, commentId, commentId)

	var up, down int64
	err := row.Scan(&up, &down)
	switch {
	case err == nil:
	case err != nil:
		return 0, 0, fmt.Errorf("row.Scan: %w", err)
	}
	return up, down, nil
}

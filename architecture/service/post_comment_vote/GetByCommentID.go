package post_comment_vote

import "fmt"

func (c *PostCommentVoteService) GetByCommentID(commentId int64) (int64, int64, error) {
	up, down, err := c.repo.GetByCommentID(commentId)
	switch {
	case err == nil:
	case err != nil:
		return 0, 0, fmt.Errorf("c.repo.GetByCommentID: %w", err)
	}
	return up, down, nil
}

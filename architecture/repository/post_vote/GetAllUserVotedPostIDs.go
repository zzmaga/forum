package post_vote

import (
	"fmt"
)

func (p *PostVoteRepo) GetAllUserVotedPostIDs(userId int64, vote int8, limit, offset int64) ([]int64, error) {
	if limit == 0 {
		limit = -1
	}

	rows, err := p.db.Query(`
	SELECT post_id FROM posts_votes
	WHERE user_id = ? AND vote = ?
	LIMIT ? OFFSET ?`, userId, vote, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("p.db.Query: %w", err)
	}

	postIDs := []int64{}
	for rows.Next() {
		var postId int64
		err = rows.Scan(&postId)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		postIDs = append(postIDs, postId)
	}

	return postIDs, nil
}

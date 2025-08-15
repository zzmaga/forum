package post_vote

import "fmt"

func (p *PostVoteRepo) GetByPostID(postId int64) (int64, int64, error) {
	row := p.db.QueryRow(`
SELECT 
	(SELECT COUNT(vote) FROM posts_votes WHERE post_id = ? AND vote == 1) as up,
    (SELECT COUNT(vote) FROM posts_votes WHERE post_id = ? AND vote == -1) as down;
	`, postId, postId)

	var up, down int64
	err := row.Scan(&up, &down)
	switch {
	case err == nil:
	case err != nil:
		return 0, 0, fmt.Errorf("row.Scan: %w", err)
	}
	return up, down, nil
}

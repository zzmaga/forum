package post_vote

import "fmt"

func (p *PostVoteService) GetByPostID(postId int64) (int64, int64, error) {
	up, down, err := p.repo.GetByPostID(postId)
	switch {
	case err == nil:
	case err != nil:
		return 0, 0, fmt.Errorf("p.repo.GetByPostID: %w", err)
	}
	return up, down, nil
}

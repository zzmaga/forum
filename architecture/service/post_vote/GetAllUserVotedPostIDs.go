package post_vote

import "fmt"

func (p *PostVoteService) GetAllUserVotedPostIDs(userId int64, vote int8, limit, offset int64) ([]int64, error) {
	postIDs, err := p.repo.GetAllUserVotedPostIDs(userId, vote, limit, offset)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("GetAllUserVotedPostIDs: %w", err)
	}
	return postIDs, nil
}

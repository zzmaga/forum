package post_vote

import (
	"forum/architecture/models"
)

type PostVoteService struct {
	repo models.IPostVoteRepo
}

func NewPostVoteService(postVote models.IPostVoteRepo) *PostVoteService {
	return &PostVoteService{postVote}
}

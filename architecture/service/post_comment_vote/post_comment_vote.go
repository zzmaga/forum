package post_comment_vote

import "forum/architecture/models"

type PostCommentVoteService struct {
	repo models.IPostCommentVoteRepo
}

func NewPostCommentVoteService(postCommentVote models.IPostCommentVoteRepo) *PostCommentVoteService {
	return &PostCommentVoteService{postCommentVote}
}

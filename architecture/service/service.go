package service

import (
	"forum/architecture/models"
	"forum/architecture/repository"
	"forum/architecture/service/category"
	"forum/architecture/service/post"
	"forum/architecture/service/post_comment"
	"forum/architecture/service/post_comment_vote"
	"forum/architecture/service/post_vote"
	"forum/architecture/service/session"
	"forum/architecture/service/user"
)

type Service struct {
	User            models.IUserService
	Post            models.IPostService
	PostVote        models.IPostVoteService
	Category        models.ICategoryService
	PostComment     models.IPostCommentService
	PostCommentVote models.IPostCommentVoteService
	Session         models.ISessionService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:            user.NewUserService(repo.User),
		Post:            post.NewPostService(repo.Post),
		PostVote:        post_vote.NewPostVoteService(repo.PostVote),
		Category:        category.NewPostCategoryService(repo.Category),
		PostComment:     post_comment.NewPostCommentService(repo.PostComment),
		PostCommentVote: post_comment_vote.NewPostCommentVoteService(repo.PostCommentVote),
		Session:         session.NewSessionService(repo.Session),
	}
}

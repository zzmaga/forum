package repository

import (
	"database/sql"

	"forum/architecture/models"
	"forum/architecture/repository/category"
	"forum/architecture/repository/post"
	"forum/architecture/repository/post_comment"
	"forum/architecture/repository/post_comment_vote"
	"forum/architecture/repository/post_vote"
	"forum/architecture/repository/session"
	"forum/architecture/repository/user"
)

type Repository struct {
	User            models.IUserRepo
	Post            models.IPostRepo
	PostVote        models.IPostVoteRepo
	Category        models.ICategoryRepo
	PostComment     models.IPostCommentRepo
	PostCommentVote models.IPostCommentVoteRepo
	Session         models.ISessionRepo
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		User:            user.NewUserRepo(db),
		Post:            post.NewPostRepo(db),
		PostVote:        post_vote.NewPostVoteRepo(db),
		Category:        category.NewCategoryRepo(db),
		PostComment:     post_comment.NewPostCommentRepo(db),
		PostCommentVote: post_comment_vote.NewPostCommentVoteRepo(db),
		Session:         session.NewSessionRepo(db),
	}
}

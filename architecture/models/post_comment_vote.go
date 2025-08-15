package models

import "time"

type PostCommentVote struct {
	Id        int64
	CommentId int64
	UserId    int64
	Vote      int8
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IPostCommentVoteRepo interface {
	Create(vote *PostCommentVote) (int64, error)
	Update(vote *PostCommentVote) error
	GetByCommentID(commentId int64) (int64, int64, error)
	GetCommentUserVote(userId, commentId int64) (*PostCommentVote, error)
	DeleteByID(id int64) error
}

type IPostCommentVoteService interface {
	Record(vote *PostCommentVote) error
	GetByCommentID(commentId int64) (int64, int64, error)
	GetCommentUserVote(userId, commentId int64) (*PostCommentVote, error)
	DeleteByID(id int64) error
}

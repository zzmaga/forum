package models

import "time"

type PostComment struct {
	Id        int64
	Content   string
	PostId    int64
	UserId    int64
	CreatedAt time.Time

	WUser     *User
	WUserVote int8  // -1 0 1
	WVoteUp   int64 // Like
	WVoteDown int64 // Dislike
}

type IPostCommentRepo interface {
	Create(comment *PostComment) (int64, error)
	// Update(comment *PostComment) error
	GetAllByPostID(postId, offset, limit int64) ([]*PostComment, error)
	GetByID(id int64) (*PostComment, error)
	DeleteByID(id int64) error
}

type IPostCommentService interface {
	Create(comment *PostComment) (int64, error)
	// Update(comment *PostComment) error
	GetAllByPostID(postId, offset, limit int64) ([]*PostComment, error)
	GetByID(id int64) (*PostComment, error)
	DeleteByID(id int64) error
}

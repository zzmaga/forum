package models

import "time"

type PostVote struct {
	Id        int64
	PostId    int64
	UserId    int64
	Vote      int8
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IPostVoteRepo interface {
	Create(vote *PostVote) (int64, error)
	Update(vote *PostVote) error
	GetByPostID(postId int64) (int64, int64, error)
	GetPostUserVote(userId, postId int64) (*PostVote, error)
	GetAllUserVotedPostIDs(userId int64, vote int8, limit, offset int64) ([]int64, error)
	DeleteByID(id int64) error
}

type IPostVoteService interface {
	Record(vote *PostVote) error
	GetByPostID(postId int64) (int64, int64, error)
	GetPostUserVote(userId, postId int64) (*PostVote, error)
	GetAllUserVotedPostIDs(userId int64, vote int8, limit, offset int64) ([]int64, error)
	DeleteByID(id int64) error
}

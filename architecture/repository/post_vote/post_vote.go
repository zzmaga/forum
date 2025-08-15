package post_vote

import (
	"database/sql"
)

type PostVoteRepo struct {
	db *sql.DB
}

func NewPostVoteRepo(db *sql.DB) *PostVoteRepo {
	return &PostVoteRepo{db}
}

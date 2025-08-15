package post_comment_vote

import "database/sql"

type PostCommentVoteRepo struct {
	db *sql.DB
}

func NewPostCommentVoteRepo(db *sql.DB) *PostCommentVoteRepo {
	return &PostCommentVoteRepo{db}
}

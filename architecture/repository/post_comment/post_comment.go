package post_comment

import "database/sql"

type PostCommentRepo struct {
	db *sql.DB
}

func NewPostCommentRepo(db *sql.DB) *PostCommentRepo {
	return &PostCommentRepo{db}
}

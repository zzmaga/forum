package post_comment

import (
	"forum/architecture/models"
)

type PostCommentService struct {
	repo models.IPostCommentRepo
}

func NewPostCommentService(repo models.IPostCommentRepo) *PostCommentService {
	return &PostCommentService{repo}
}

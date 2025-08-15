package post

import "forum/architecture/models"

type PostService struct {
	repo models.IPostRepo
}

func NewPostService(repo models.IPostRepo) *PostService {
	return &PostService{repo}
}

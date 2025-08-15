package service

import "forum/architecture/models"

type Service struct {
	User models.IUserService
	// Post models.IPostService
}

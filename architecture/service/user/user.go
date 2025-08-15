package user

import "forum/architecture/models"

type UserService struct {
	repo models.IUserRepo
}

func NewUserService(repo models.IUserRepo) *UserService {
	return &UserService{repo}
}

func (u *UserService) GetByEmail(email string) (*models.User, error) {
	return u.repo.GetByEmail(email)
}

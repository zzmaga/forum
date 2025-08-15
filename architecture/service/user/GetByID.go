package user

import (
	"errors"
	"fmt"

	"forum/architecture/models"
	"forum/architecture/repository/user"
)

func (u *UserService) GetByID(id int64) (*models.User, error) {
	usr, err := u.repo.GetByID(id)
	switch {
	case err == nil:
		return usr, nil
	case errors.Is(err, user.ErrNotFound):
		return nil, ErrNotFound
	}
	return nil, fmt.Errorf("u.repo.GetByID: %w", err)
}

package user

import (
	"errors"
	"fmt"
	"strings"

	"forum/architecture/models"

	"forum/architecture/repository/user"
)

func (u *UserService) GetByNicknameOrEmail(field string) (*models.User, error) {
	switch {
	case strings.Contains(field, "@"):
		if err := (&models.User{Email: field}).ValidateEmail(); err != nil {
			return nil, ErrInvalidEmail
		}
		usr, err := u.repo.GetByEmail(field)
		switch {
		case err == nil:
			return usr, err
		case errors.Is(err, user.ErrNotFound):
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("u.repo.GetByEmail: %w", err)
	default:
		if err := (&models.User{Nickname: field}).ValidateNickname(); err != nil {
			return nil, ErrInvalidNickname
		}
		usr, err := u.repo.GetByNickname(field)
		switch {
		case err == nil:
			return usr, err
		case errors.Is(err, user.ErrNotFound):
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("u.repo.GetByNickname: %w", err)
	}
}

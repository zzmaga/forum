package user

import (
	"fmt"
	"strings"

	"forum/architecture/models"
)

func (u *UserRepo) Create(user *models.User) (int64, error) {
	strCreatedAt := user.CreatedAt.Format(models.TimeFormat)
	row := u.db.QueryRow(`
INSERT INTO users (nickname, email, password, created_at) VALUES
(?, ?, ?, ?) RETURNING id`, user.Nickname, user.Email, user.Password, strCreatedAt)

	err := row.Scan(&user.Id)
	switch {
	case err == nil:
		return user.Id, nil
	case strings.HasPrefix(err.Error(), "UNIQUE constraint failed"):
		switch {
		case strings.Contains(err.Error(), "nickname"):
			return -1, ErrExistNickname
		case strings.Contains(err.Error(), "email"):
			return -1, ErrExistEmail
		}
	case strings.HasPrefix(err.Error(), "CHECK constraint failed"):
		switch {
		case strings.Contains(err.Error(), "LENGTH(nickname)"):
			return -1, ErrWrongLengthNickname
		case strings.Contains(err.Error(), "LENGTH(email)"):
			return -1, ErrWrongLengthEmail
		}
	}
	return -1, fmt.Errorf("row.Scan: %w", err)
}

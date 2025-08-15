package user

import (
	"fmt"
	"strings"
	"time"

	"forum/architecture/models"
)

func (u *UserRepo) GetByID(id int64) (*models.User, error) {
	row := u.db.QueryRow(`
SELECT id, nickname, email, created_at FROM users
WHERE id = ?`, id)
	user := &models.User{}
	strCreatedAt := ""
	err := row.Scan(&user.Id, &user.Nickname, &user.Email, &strCreatedAt)
	switch {
	case err == nil:
		timeCreatedAt, err := time.ParseInLocation(models.TimeFormat, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		user.CreatedAt = timeCreatedAt
		return user, nil
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
}

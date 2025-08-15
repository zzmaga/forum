package session

import (
	"fmt"
	"strings"

	"forum/architecture/models"
)

func (s *SessionRepo) Create(session *models.Session) (int64, error) {
	strExpiredAt := session.ExpiredAt.Format(models.TimeFormat)
	row := s.db.QueryRow(`
INSERT INTO sessions (uuid, expired_at, user_id) VALUES
(?, ?, ?) RETURNING id`, session.Uuid, strExpiredAt, session.UserId)

	err := row.Scan(&session.Id)
	switch {
	case err == nil:
		return session.Id, nil
	case strings.HasPrefix(err.Error(), "UNIQUE constraint failed"):
		return -1, ErrSessionExists
	}
	return -1, fmt.Errorf("row.Scan: %w", err)
}

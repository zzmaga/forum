package session

import (
	"fmt"
	"strings"
	"time"

	"forum/architecture/models"
)

func (s *SessionRepo) GetByUuid(uuid string) (*models.Session, error) {
	row := s.db.QueryRow(`
SELECT id, uuid, expired_at, user_id FROM sessions
WHERE uuid = ?`, uuid)
	session := &models.Session{}
	strExpiredAt := ""

	err := row.Scan(&session.Id, &session.Uuid, &strExpiredAt, &session.UserId)

	switch {
	case err == nil:
		timeExpiredAt, err := time.ParseInLocation(models.TimeFormat, strExpiredAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		session.ExpiredAt = timeExpiredAt
		return session, nil
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
}

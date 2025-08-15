package session

import (
	"fmt"

	"forum/architecture/models"
)

func (s *SessionRepo) UpdateByUserId(userId int64, session *models.Session) error {
	strExpiredAt := session.ExpiredAt.Format(models.TimeFormat)
	row := s.db.QueryRow(`
UPDATE sessions 
SET uuid = ?, expired_at = ?
WHERE user_id = ?
RETURNING id`, session.Uuid, strExpiredAt, session.UserId)

	err := row.Scan(&session.Id)
	switch {
	case err == nil:
		return nil
	}
	return fmt.Errorf("row.Scan: %w", err)
}

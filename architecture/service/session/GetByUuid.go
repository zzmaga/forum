package session

import (
	"errors"
	"fmt"
	"time"

	"forum/architecture/models"
	rsession "forum/architecture/repository/session"
)

func (s *SessionService) GetByUuid(uuid string) (*models.Session, error) {
	session, err := s.repo.GetByUuid(uuid)
	switch {
	case err == nil:
		expiredInSec := time.Until(session.ExpiredAt).Seconds()
		if expiredInSec <= 0 {
			return nil, ErrExpired
		}
		return session, nil
	case errors.Is(err, rsession.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("s.repo.GetByUuid: %w", err)
	}
}

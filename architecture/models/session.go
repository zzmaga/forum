package models

import "time"

type Session struct {
	Id        int64
	Uuid      string
	ExpiredAt time.Time
	UserId    int64
}

type ISessionRepo interface {
	Create(session *Session) (int64, error)
	Delete(id int64) error
	GetByUuid(uuid string) (*Session, error)
	UpdateByUserId(userId int64, session *Session) error
}

type ISessionService interface {
	Record(userId int64) (*Session, error)
	Delete(id int64) error
	GetByUuid(uuid string) (*Session, error)
}

package models

import (
	"time"
)

// User -
type User struct {
	Id        int64
	Nickname  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type IUserService interface {
	Create(user *User) (int64, error)
	Update(user *User) error
	DeleteByID(id int64) error

	GetByID(id int64) (*User, error)
	GetByNicknameOrEmail(field string) (*User, error)
	// GetAll(from, offset int64) error
}

type IUserRepo interface {
	Create(user *User) (int64, error)
	Update(user *User) error
	DeleteByID(id int64) error

	GetByID(id int64) (*User, error)
	GetByNickname(nickname string) (*User, error)
	GetByEmail(email string) (*User, error)
	// GetAll(from, offset int64) error
}

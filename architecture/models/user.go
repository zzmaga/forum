package models

type User struct {
	Id       int64
	Username string
	Email    string
	Password string
}

type IUserService interface {
	Create(user *User) (int64, error)
	Update(user *User) error
	DeleteByID(id int64) error

	GetByID(id int64) (*User, error)
	GetByUsernameOrEmail(field string) (*User, error)

	// GetAll(from, offset int64) error
}

type IUserRepo interface {
	Create(user *User) (int64, error)
	Update(user *User) error
	DeleteByID(id int64) error

	GetByID(id int64) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
}

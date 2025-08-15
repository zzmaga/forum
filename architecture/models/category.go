package models

import "time"

type Category struct {
	Id        int64
	Name      string
	CreatedAt time.Time
}

type ICategoryRepo interface {
	Create(category *Category) (int64, error)
	AddToPost(categoryId, postId int64) (int64, error)
	Update(category *Category) error
	GetByID(id int64) (*Category, error)
	GetByName(name string) (*Category, error)
	GetByNames(names []string) ([]*Category, error)
	GetByPostID(postId int64) ([]*Category, error)
	GetPostIDsContainedCatIDs(ids []int64, offset, limit int64) ([]int64, error)
	GetAll(offset, limit int64) ([]*Category, error)
	DeleteByPostID(postId int64) error
	DeleteByID(id int64) error
}

type ICategoryService interface {
	AddToPostByNames(names []string, postId int64) error
	GetByPostID(postId int64) ([]*Category, error)
	GetByNames(names []string) ([]*Category, error)
	GetAll(offset, limit int64) ([]*Category, error)
	GetPostIDsContainedCatIDs(ids []int64, offset, limit int64) ([]int64, error)
	DeleteByPostID(postId int64) error
	DeleteFromPost(id int64) error
}

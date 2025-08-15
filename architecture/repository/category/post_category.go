package category

import (
	"database/sql"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db}
}

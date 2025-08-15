package category

import "errors"

var (
	ErrExistName       = errors.New("category with this name exists")
	ErrCheckLengthName = errors.New("category name length is wrong")
	ErrNotFound        = errors.New("category not found")
)

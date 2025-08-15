package category

import (
	"fmt"
)

func (c *CategoryRepo) AddToPost(categoryId, postId int64) (int64, error) {
	row := c.db.QueryRow(`
	INSERT INTO posts_categories (post_id, category_id) VALUES
	(?, ?) RETURNING id`, postId, categoryId)

	id := int64(-1)
	err := row.Scan(&id)
	switch {
	case err == nil:
		return id, nil
	}
	return -1, fmt.Errorf("row.Scan: %w", err)
}

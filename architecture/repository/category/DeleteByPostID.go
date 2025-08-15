package category

import "fmt"

func (c *CategoryRepo) DeleteByPostID(postId int64) error {
	_, err := c.db.Exec(`DELETE FROM posts_categories
WHERE post_id = ?`, postId)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("c.db.Exec: %w", err)
	}
	return nil
}

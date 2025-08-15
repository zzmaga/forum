package post_comment

import "fmt"

func (c *PostCommentRepo) DeleteByID(id int64) error {
	_, err := c.db.Exec("DELETE FROM posts_comments WHERE id = ?", id)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("p.db.Exec: %w", err)
	}
	return nil
}

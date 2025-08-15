package post_comment

import "fmt"

func (c *PostCommentService) DeleteByID(id int64) error {
	err := c.repo.DeleteByID(id)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("c.repo.DeleteByID: %w", err)
	}
	return nil
}

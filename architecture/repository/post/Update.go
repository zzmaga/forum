package post

import (
	"fmt"

	"forum/architecture/models"
)

func (p *PostRepo) Update(post *models.Post) error {
	strUpdatedAt := post.UpdatedAt.Format(models.TimeFormat)
	_, err := p.db.Exec(`UPDATE posts
SET title = ?, content = ?, updated_at = ?
WHERE id = ?`, post.Title, post.Content, strUpdatedAt, post.Id)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("p.db.Exec: %w", err)
	}
	return nil
}

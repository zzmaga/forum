package post

import "fmt"

func (p *PostRepo) DeleteByID(id int64) error {
	_, err := p.db.Exec("DELETE FROM posts WHERE id = ?", id)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("p.db.Exec: %w", err)
	}
	return nil
}

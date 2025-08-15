package post

import "fmt"

func (q *PostService) DeleteByID(id int64) error {
	err := q.repo.DeleteByID(id)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("q.repo.DeleteByID: %w", err)
	}
	return nil
}

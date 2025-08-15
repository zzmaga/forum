package category

import (
	"errors"
	"fmt"
	"time"

	"forum/architecture/models"

	rcategory "forum/architecture/repository/category"
)

func (c *CategoryService) AddToPostByNames(names []string, postId int64) error {
	if len(names) == 0 {
		return nil
	} else if len(names) > models.MaxCategoryLimitForPost {
		return ErrCategoryLimitForPost
	}

	var ids []int64 = make([]int64, len(names))
	for i, name := range names {
		cat := &models.Category{Name: name, CreatedAt: time.Now()}
		id, err := c.repo.Create(cat)
		switch {
		case err == nil:
			ids[i] = id
			continue
		case errors.Is(err, rcategory.ErrExistName):
		case errors.Is(err, rcategory.ErrCheckLengthName):
			return ErrCheckLengthName
		default:
			return fmt.Errorf("c.repo.Create: %w", err)
		}

		cat, err = c.repo.GetByName(name)
		switch {
		case err == nil:
			ids[i] = cat.Id
			continue
		default:
			return fmt.Errorf("c.repo.GetByName: %w", err)
		}
	}

	for _, id := range ids {
		_, err := c.repo.AddToPost(id, postId)
		switch {
		case err == nil:
			continue
		default:
			return fmt.Errorf("c.repo.AddToPost: %w", err)
		}
	}
	return nil
}

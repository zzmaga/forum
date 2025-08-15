package category

import (
	"fmt"
)

func (c *CategoryService) GetPostIDsContainedCatIDs(ids []int64, offset, limit int64) ([]int64, error) {
	cats, err := c.repo.GetPostIDsContainedCatIDs(ids, offset, limit)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("GetPostIDsContainedCatIDs: %w", err)
	}
	return cats, nil
}

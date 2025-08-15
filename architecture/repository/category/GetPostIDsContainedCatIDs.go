package category

import (
	"fmt"
	"strings"
)

func (c *CategoryRepo) GetPostIDsContainedCatIDs(ids []int64, offset, limit int64) ([]int64, error) {
	strIDs := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
	preQuery := fmt.Sprintf(`SELECT post_id, COUNT(category_id) as cat from posts_categories
WHERE category_id IN (%s)
GROUP BY post_id
HAVING cat >= %d
LIMIT ? OFFSET ?`, strIDs, len(ids))

	rows, err := c.db.Query(preQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}

	postIDs := []int64{}
	for rows.Next() {
		var postId, a int64
		err = rows.Scan(&postId, &a)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		postIDs = append(postIDs, postId)
	}
	return postIDs, nil
}

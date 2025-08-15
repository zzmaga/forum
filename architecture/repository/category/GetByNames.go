package category

import (
	"fmt"
	"strings"
	"time"

	"forum/architecture/models"
)

func (c *CategoryRepo) GetByNames(names []string) ([]*models.Category, error) {
	if len(names) == 0 {
		return nil, nil
	}

	iNames := make([]interface{}, len(names))
	for i, v := range names {
		iNames[i] = v
	}

	strQuery := fmt.Sprintf(`SELECT id, name, created_at FROM categories
WHERE name IN (%v)`, `?`+strings.Repeat(",?", len(iNames)-1))
	rows, err := c.db.Query(strQuery, iNames...)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}

	cats := []*models.Category{}
	for rows.Next() {
		var strCreatedAt string
		cat := &models.Category{}
		err = rows.Scan(&cat.Id, &cat.Name, &strCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		timeCreatedAt, err := time.ParseInLocation(models.TimeFormat, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		cat.CreatedAt = timeCreatedAt

		cats = append(cats, cat)
	}
	return cats, nil
}

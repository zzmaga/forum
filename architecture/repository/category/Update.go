package category

import (
	"database/sql"
	"forum/architecture/models"
)

// Update обновляет существующую категорию
func (c *CategoryRepo) Update(category *models.Category) error {
	query := "UPDATE categories SET name = ? WHERE id = ?"

	result, err := c.db.Exec(query, category.Name, category.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

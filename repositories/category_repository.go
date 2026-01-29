package repositories

import (
	"cashier-api/models"
	"database/sql"
	"errors"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name, description FROM categories ORDER BY id"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"
	row := r.db.QueryRow(query, id)

	var c models.Category
	err := row.Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &c, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	return r.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
}

func (r *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	result, err := r.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("category not found")
	}

	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	// Check if category is being used by any products
	var productCount int
	checkQuery := "SELECT COUNT(*) FROM products WHERE category_id = $1"
	err := r.db.QueryRow(checkQuery, id).Scan(&productCount)
	if err != nil {
		return err
	}

	if productCount > 0 {
		return errors.New("cannot delete category that has products")
	}

	query := "DELETE FROM categories WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("category not found")
	}

	return nil
}

package repositories

import (
	"cashier-api/models"
	"database/sql"
	"errors"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll - Get all products WITHOUT category info
func (r *ProductRepository) GetAll() ([]models.ProductList, error) {
	query := "SELECT id, name, price, stock FROM products ORDER BY id"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.ProductList
	for rows.Next() {
		var p models.ProductList
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetByID - Get product by ID WITH category name (JOIN)
func (r *ProductRepository) GetByID(id int) (*models.ProductDetail, error) {
	query := `
        SELECT p.id, p.name, p.price, p.stock, p.category_id, 
               c.name as category_name
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
        WHERE p.id = $1
    `
	row := r.db.QueryRow(query, id)

	var product models.ProductDetail

	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock,
		&product.CategoryID, &product.CategoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

// Create - Create new product
func (r *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	return r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
}

// Update - Update product
func (r *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := r.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

// Delete - Delete product
func (r *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

// CheckCategoryExists - Helper to validate category_id
func (r *ProductRepository) CheckCategoryExists(categoryID int) (bool, error) {
	query := "SELECT COUNT(*) FROM categories WHERE id = $1"
	var count int
	err := r.db.QueryRow(query, categoryID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

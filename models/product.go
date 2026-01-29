package models

import "errors"

// Product - Basic product structure for create/update
type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID int    `json:"category_id"`
}

// ProductList - For GET /api/products response (NO category)
type ProductList struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

// ProductDetail - For GET /api/products/{id} response (WITH category_name)
type ProductDetail struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Stock        int    `json:"stock"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"` // Only in detail
}

// Validation errors
var (
	ErrInvalidID         = errors.New("invalid ID")
	ErrNameRequired      = errors.New("name is required")
	ErrInvalidPrice      = errors.New("price must be greater than 0")
	ErrInvalidStock      = errors.New("stock cannot be negative")
	ErrInvalidCategoryID = errors.New("invalid category ID")
	ErrCategoryNotFound  = errors.New("category not found")
)

package models

import "errors"

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

// Common validation errors
var (
	ErrInvalidID    = errors.New("invalid ID")
	ErrNameRequired = errors.New("name is required")
	ErrInvalidPrice = errors.New("price must be greater than 0")
	ErrInvalidStock = errors.New("stock cannot be negative")
)

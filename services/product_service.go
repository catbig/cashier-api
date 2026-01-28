package services

import (
	"cashier-api/models"
	"cashier-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	if id <= 0 {
		return nil, models.ErrInvalidID
	}
	return s.repo.GetByID(id)
}

func (s *ProductService) Create(product *models.Product) error {
	if product.Name == "" {
		return models.ErrNameRequired
	}
	if product.Price <= 0 {
		return models.ErrInvalidPrice
	}
	if product.Stock < 0 {
		return models.ErrInvalidStock
	}
	return s.repo.Create(product)
}

func (s *ProductService) Update(product *models.Product) error {
	if product.ID <= 0 {
		return models.ErrInvalidID
	}
	if product.Name == "" {
		return models.ErrNameRequired
	}
	if product.Price <= 0 {
		return models.ErrInvalidPrice
	}
	if product.Stock < 0 {
		return models.ErrInvalidStock
	}
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	if id <= 0 {
		return models.ErrInvalidID
	}
	return s.repo.Delete(id)
}

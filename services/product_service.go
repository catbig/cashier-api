package services

import (
	"cashier-api/models"
	"cashier-api/repositories"
)

type ProductService struct {
	productRepo  *repositories.ProductRepository
	categoryRepo *repositories.CategoryRepository
}

func NewProductService(productRepo *repositories.ProductRepository, categoryRepo *repositories.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductService) GetAll() ([]models.ProductList, error) {
	return s.productRepo.GetAll()
}

func (s *ProductService) GetByID(id int) (*models.ProductDetail, error) {
	if id <= 0 {
		return nil, models.ErrInvalidID
	}
	return s.productRepo.GetByID(id)
}

func (s *ProductService) Create(product *models.Product) error {
	// Basic validation
	if product.Name == "" {
		return models.ErrNameRequired
	}
	if product.Price <= 0 {
		return models.ErrInvalidPrice
	}
	if product.Stock < 0 {
		return models.ErrInvalidStock
	}
	if product.CategoryID <= 0 {
		return models.ErrInvalidCategoryID
	}

	// Validate category exists
	categoryExists, err := s.productRepo.CheckCategoryExists(product.CategoryID)
	if err != nil {
		return err
	}
	if !categoryExists {
		return models.ErrCategoryNotFound
	}

	return s.productRepo.Create(product)
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
	if product.CategoryID <= 0 {
		return models.ErrInvalidCategoryID
	}

	// Validate category exists
	categoryExists, err := s.productRepo.CheckCategoryExists(product.CategoryID)
	if err != nil {
		return err
	}
	if !categoryExists {
		return models.ErrCategoryNotFound
	}

	return s.productRepo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	if id <= 0 {
		return models.ErrInvalidID
	}
	return s.productRepo.Delete(id)
}

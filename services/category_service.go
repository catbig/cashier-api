package services

import (
	"cashier-api/models"
	"cashier-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	if id <= 0 {
		return nil, models.ErrInvalidID
	}
	return s.repo.GetByID(id)
}

func (s *CategoryService) Create(category *models.Category) error {
	if category.Name == "" {
		return models.ErrNameRequired
	}
	return s.repo.Create(category)
}

func (s *CategoryService) Update(category *models.Category) error {
	if category.ID <= 0 {
		return models.ErrInvalidID
	}
	if category.Name == "" {
		return models.ErrNameRequired
	}
	return s.repo.Update(category)
}

func (s *CategoryService) Delete(id int) error {
	if id <= 0 {
		return models.ErrInvalidID
	}
	return s.repo.Delete(id)
}

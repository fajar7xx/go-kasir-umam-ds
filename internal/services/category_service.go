package services

import (
	"fajar7xx/go-kasir-umam-ds/internal/repositories"
	"fajar7xx/go-kasir-umam-ds/models"
)

type CategoryServiceInterface interface {
	GetAll() ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	Create(category *models.Category) error
	Update(id int, category *models.Category) (*models.Category, error)
	Delete(id int) error
}

type CategoryService struct {
	categoryRepo repositories.CategoryRepositoryInterface
}

func NewCategoryService(categoryRepo repositories.CategoryRepositoryInterface) CategoryServiceInterface {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (serv *CategoryService) GetAll() ([]models.Category, error) {
	return serv.categoryRepo.GetAll()
}

func (serv *CategoryService) GetByID(id int) (*models.Category, error) {
	return serv.categoryRepo.GetByID(id)
}

func (serv *CategoryService) Create(category *models.Category) error {
	return serv.categoryRepo.Create(category)
}

func (serv *CategoryService) Update(id int, category *models.Category) (*models.Category, error) {
	err := serv.categoryRepo.Update(id, category)
	if err != nil {
		return nil, err
	}

	updatedCategory, err := serv.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return updatedCategory, nil
}

func (serv *CategoryService) Delete(id int) error {
	return serv.categoryRepo.Delete(id)
}

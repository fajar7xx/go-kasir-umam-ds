package services

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/internal/repositories"
	"fajar7xx/go-kasir-umam-ds/models"
)

type CategoryService interface {
	GetAll(ctx context.Context, name string) ([]models.CategoryResponse, error)
	GetByID(ctx context.Context, id int) (*models.CategoryResponse, error)
	Create(ctx context.Context, category *models.Category) (*models.CategoryResponse, error)
	Update(ctx context.Context, id int, category *models.Category) (*models.CategoryResponse, error)
	Delete(ctx context.Context, id int) error
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (serv *categoryService) GetAll(ctx context.Context, name string) ([]models.CategoryResponse, error) {
	return serv.categoryRepo.GetAll(ctx, name)
}

func (serv *categoryService) GetByID(ctx context.Context, id int) (*models.CategoryResponse, error) {
	return serv.categoryRepo.GetByID(ctx, id)
}

func (serv *categoryService) Create(ctx context.Context, category *models.Category) (*models.CategoryResponse, error) {
	return serv.categoryRepo.Create(ctx, category)
}

func (serv *categoryService) Update(ctx context.Context, id int, category *models.Category) (*models.CategoryResponse, error) {
	err := serv.categoryRepo.Update(ctx, id, category)
	if err != nil {
		return nil, err
	}

	updatedCategory, err := serv.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedCategory, nil
}

func (serv *categoryService) Delete(ctx context.Context, id int) error {
	return serv.categoryRepo.Delete(ctx, id)
}

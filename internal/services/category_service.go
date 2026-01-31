package services

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/internal/repositories"
	"fajar7xx/go-kasir-umam-ds/models"
)

type CategoryServiceInterface interface {
	GetAll(ctx context.Context) ([]models.CategoryResponse, error)
	GetByID(ctx context.Context, id int) (*models.CategoryResponse, error)
	Create(ctx context.Context, category *models.Category) (*models.CategoryResponse, error)
	Update(ctx context.Context, id int, category *models.Category) (*models.CategoryResponse, error)
	Delete(ctx context.Context, id int) error
}

type CategoryService struct {
	categoryRepo repositories.CategoryRepositoryInterface
}

func NewCategoryService(categoryRepo repositories.CategoryRepositoryInterface) CategoryServiceInterface {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (serv *CategoryService) GetAll(ctx context.Context) ([]models.CategoryResponse, error) {
	return serv.categoryRepo.GetAll(ctx)
}

func (serv *CategoryService) GetByID(ctx context.Context, id int) (*models.CategoryResponse, error) {
	return serv.categoryRepo.GetByID(ctx, id)
}

func (serv *CategoryService) Create(ctx context.Context, category *models.Category) (*models.CategoryResponse, error) {
	return serv.categoryRepo.Create(ctx, category)
}

func (serv *CategoryService) Update(ctx context.Context, id int, category *models.Category) (*models.CategoryResponse, error) {
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

func (serv *CategoryService) Delete(ctx context.Context, id int) error {
	return serv.categoryRepo.Delete(ctx, id)
}

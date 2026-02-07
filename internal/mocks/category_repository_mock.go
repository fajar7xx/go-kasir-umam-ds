package mocks

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/models"

	"github.com/stretchr/testify/mock"
)

type CategoryRepositoryMock struct {
	mock.Mock
}

func (m *CategoryRepositoryMock) GetAll(ctx context.Context, name string) ([]models.CategoryResponse, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CategoryResponse), args.Error(1)
}

func (m *CategoryRepositoryMock) GetByID(ctx context.Context, id int) (*models.CategoryResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CategoryResponse), args.Error(1)
}

func (m *CategoryRepositoryMock) Create(ctx context.Context, category *models.Category) (*models.CategoryResponse, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CategoryResponse), args.Error(1)
}

func (m *CategoryRepositoryMock) Update(ctx context.Context, id int, category *models.Category) error {
	args := m.Called(ctx, id, category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

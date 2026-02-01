package mocks

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/models"

	"github.com/stretchr/testify/mock"
)

type CategoryServiceMock struct {
	mock.Mock
}

func (m *CategoryServiceMock) GetAll(ctx context.Context) ([]models.CategoryResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CategoryResponse), args.Error(1)
}

func (m *CategoryServiceMock) GetByID(ctx context.Context, id int) (*models.CategoryResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CategoryResponse), args.Error(1)
}

func (m *CategoryServiceMock) Create(ctx context.Context, category *models.Category) (*models.CategoryResponse, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CategoryResponse), args.Error(1)
}

func (m *CategoryServiceMock) Update(ctx context.Context, id int, category *models.Category) (*models.CategoryResponse, error) {
	args := m.Called(ctx, id, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CategoryResponse), args.Error(1)
}

func (m *CategoryServiceMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

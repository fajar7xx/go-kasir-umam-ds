package mocks

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/models"

	"github.com/stretchr/testify/mock"
)

type ProductServiceMock struct {
	mock.Mock
}

func (m *ProductServiceMock) GetAll(ctx context.Context, name string) ([]models.ProductResponse, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ProductResponse), args.Error(1)
}

func (m *ProductServiceMock) GetByID(ctx context.Context, id int) (*models.ProductResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductResponse), args.Error(1)
}

func (m *ProductServiceMock) Create(ctx context.Context, product *models.Product) (*models.ProductResponse, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductResponse), args.Error(1)
}

func (m *ProductServiceMock) Update(ctx context.Context, id int, product *models.Product) (*models.ProductResponse, error) {
	args := m.Called(ctx, id, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductResponse), args.Error(1)
}

func (m *ProductServiceMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

package mocks

import (
	"fajar7xx/go-kasir-umam-ds/models"
	"github.com/stretchr/testify/mock"
)

type ProductServiceMock struct {
	mock.Mock
}

func (m *ProductServiceMock) GetAll() ([]models.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *ProductServiceMock) GetByID(id int) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductServiceMock) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *ProductServiceMock) Update(id int, product *models.Product) (*models.Product, error) {
	args := m.Called(id, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductServiceMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

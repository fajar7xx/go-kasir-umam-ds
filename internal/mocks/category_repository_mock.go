package mocks

import (
	"fajar7xx/go-kasir-umam-ds/models"
	"github.com/stretchr/testify/mock"
)

type CategoryRepositoryMock struct {
	mock.Mock
}

func (m *CategoryRepositoryMock) GetAll() ([]models.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) GetByID(id int) (*models.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) Create(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) Update(id int, category *models.Category) error {
	args := m.Called(id, category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

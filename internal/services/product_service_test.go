package services

import (
	"errors"
	"fajar7xx/go-kasir-umam-ds/internal/mocks"
	"fajar7xx/go-kasir-umam-ds/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProductService_GetAll(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	now := time.Now()
	expectedProducts := []models.Product{
		{ID: 1, Name: "Nasi Goreng", CreatedAt: now},
		{ID: 2, Name: "Es Teh", CreatedAt: now},
	}

	mockRepo.On("GetAll").Return(expectedProducts, nil)

	products, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Nasi Goreng", products[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetAll_Error(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	mockRepo.On("GetAll").Return(nil, errors.New("database error"))

	products, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, products)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetByID(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	now := time.Now()
	expectedProduct := &models.Product{ID: 1, Name: "Nasi Goreng", CreatedAt: now}

	mockRepo.On("GetByID", 1).Return(expectedProduct, nil)

	product, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, product.ID)
	assert.Equal(t, "Nasi Goreng", product.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Create(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	product := &models.Product{Name: "Nasi Goreng"}

	mockRepo.On("Create", product).Return(nil)

	err := service.Create(product)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	id := 1
	product := &models.Product{Name: "Nasi Goreng Updated"}
	updatedProduct := &models.Product{ID: 1, Name: "Nasi Goreng Updated"}

	// Expect Update to be called
	mockRepo.On("Update", id, product).Return(nil)
	// Expect GetByID to be called after Update
	mockRepo.On("GetByID", id).Return(updatedProduct, nil)

	result, err := service.Update(id, product)

	assert.NoError(t, err)
	assert.Equal(t, "Nasi Goreng Updated", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update_Error(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	id := 1
	product := &models.Product{Name: "Nasi Goreng Updated"}

	mockRepo.On("Update", id, product).Return(errors.New("update failed"))

	result, err := service.Update(id, product)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Delete(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	id := 1
	mockRepo.On("Delete", id).Return(nil)

	err := service.Delete(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

package services

import (
	"context"
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
	expectedProducts := []models.ProductResponse{
		{ID: 1, Name: "Nasi Goreng", CreatedAt: now},
		{ID: 2, Name: "Es Teh", CreatedAt: now},
	}

	mockRepo.On("GetAll", context.Background()).Return(expectedProducts, nil)

	products, err := service.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Nasi Goreng", products[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetAll_Error(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	mockRepo.On("GetAll", context.Background()).Return(nil, errors.New("database error"))

	products, err := service.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, products)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetByID(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	now := time.Now()
	expectedProduct := &models.ProductResponse{ID: 1, Name: "Nasi Goreng", CreatedAt: now}

	mockRepo.On("GetByID", context.Background(), 1).Return(expectedProduct, nil)

	product, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, product.ID)
	assert.Equal(t, "Nasi Goreng", product.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Create(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	product := &models.Product{Name: "Nasi Goreng"}
	expectedResponse := &models.ProductResponse{ID: 1, Name: "Nasi Goreng"}

	mockRepo.On("Create", context.Background(), product).Return(expectedResponse, nil)

	result, err := service.Create(context.Background(), product)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Nasi Goreng", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	id := 1
	product := &models.Product{Name: "Nasi Goreng Updated"}
	updatedProduct := &models.ProductResponse{ID: 1, Name: "Nasi Goreng Updated"}

	// Expect Update to be called
	mockRepo.On("Update", context.Background(), id, product).Return(nil)
	// Expect GetByID to be called after Update
	mockRepo.On("GetByID", context.Background(), id).Return(updatedProduct, nil)

	result, err := service.Update(context.Background(), id, product)

	assert.NoError(t, err)
	assert.Equal(t, "Nasi Goreng Updated", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update_Error(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	id := 1
	product := &models.Product{Name: "Nasi Goreng Updated"}

	mockRepo.On("Update", context.Background(), id, product).Return(errors.New("update failed"))

	result, err := service.Update(context.Background(), id, product)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Delete(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	id := 1
	mockRepo.On("Delete", context.Background(), id).Return(nil)

	err := service.Delete(context.Background(), id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update_GetByIDErrorAfterUpdate(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	service := NewProductService(mockRepo)

	id := 1
	desc := "Product Updated"
	product := &models.Product{
		Name:        "Product Updated",
		Description: &desc,
		Price:       15000,
		Stock:       10,
		CategoryID:  1,
	}

	mockRepo.On("Update", context.Background(), id, product).Return(nil)
	mockRepo.On("GetByID", context.Background(), id).Return(nil, errors.New("database connection lost"))

	result, err := service.Update(context.Background(), id, product)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database connection lost", err.Error())
	mockRepo.AssertExpectations(t)
}

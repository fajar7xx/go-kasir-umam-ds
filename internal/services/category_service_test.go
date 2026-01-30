package services

import (
	"errors"
	"fajar7xx/go-kasir-umam-ds/internal/mocks"
	"fajar7xx/go-kasir-umam-ds/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCategoryService_GetAll(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	now := time.Now()
	expectedCategories := []models.Category{
		{ID: 1, Name: "Food", CreatedAt: now},
		{ID: 2, Name: "Beverage", CreatedAt: now},
	}

	mockRepo.On("GetAll").Return(expectedCategories, nil)

	categories, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Food", categories[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetAll_Error(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	mockRepo.On("GetAll").Return(nil, errors.New("database error"))

	categories, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, categories)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetByID(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	now := time.Now()
	expectedCategory := &models.Category{ID: 1, Name: "Food", CreatedAt: now}

	mockRepo.On("GetByID", 1).Return(expectedCategory, nil)

	category, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, category.ID)
	assert.Equal(t, "Food", category.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Create(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	category := &models.Category{Name: "Food"}

	mockRepo.On("Create", category).Return(nil)

	err := service.Create(category)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	id := 1
	category := &models.Category{Name: "Food Updated"}
	updatedCategory := &models.Category{ID: 1, Name: "Food Updated"}

	// Expect Update to be called
	mockRepo.On("Update", id, category).Return(nil)
	// Expect GetByID to be called after Update
	mockRepo.On("GetByID", id).Return(updatedCategory, nil)

	result, err := service.Update(id, category)

	assert.NoError(t, err)
	assert.Equal(t, "Food Updated", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update_Error(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	id := 1
	category := &models.Category{Name: "Food Updated"}

	mockRepo.On("Update", id, category).Return(errors.New("update failed"))

	result, err := service.Update(id, category)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Delete(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	id := 1
	mockRepo.On("Delete", id).Return(nil)

	err := service.Delete(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

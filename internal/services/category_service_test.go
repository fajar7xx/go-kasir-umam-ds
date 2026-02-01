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

func TestCategoryService_GetAll(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	now := time.Now()
	expectedCategories := []models.CategoryResponse{
		{ID: 1, Name: "Food", CreatedAt: now},
		{ID: 2, Name: "Beverage", CreatedAt: now},
	}

	mockRepo.On("GetAll", context.Background()).Return(expectedCategories, nil)

	categories, err := service.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Food", categories[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetAll_Error(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	mockRepo.On("GetAll", context.Background()).Return(nil, errors.New("database error"))

	categories, err := service.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, categories)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetByID(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	now := time.Now()
	expectedCategory := &models.CategoryResponse{ID: 1, Name: "Food", CreatedAt: now}

	mockRepo.On("GetByID", context.Background(), 1).Return(expectedCategory, nil)

	category, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, category.ID)
	assert.Equal(t, "Food", category.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Create(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	category := &models.Category{Name: "Food"}
	expectedResponse := &models.CategoryResponse{ID: 1, Name: "Food"}

	mockRepo.On("Create", context.Background(), category).Return(expectedResponse, nil)

	result, err := service.Create(context.Background(), category)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Food", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	id := 1
	category := &models.Category{Name: "Food Updated"}
	updatedCategory := &models.CategoryResponse{ID: 1, Name: "Food Updated"}

	// Expect Update to be called
	mockRepo.On("Update", context.Background(), id, category).Return(nil)
	// Expect GetByID to be called after Update
	mockRepo.On("GetByID", context.Background(), id).Return(updatedCategory, nil)

	result, err := service.Update(context.Background(), id, category)

	assert.NoError(t, err)
	assert.Equal(t, "Food Updated", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update_Error(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	id := 1
	category := &models.Category{Name: "Food Updated"}

	mockRepo.On("Update", context.Background(), id, category).Return(errors.New("update failed"))

	result, err := service.Update(context.Background(), id, category)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Delete(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	id := 1
	mockRepo.On("Delete", context.Background(), id).Return(nil)

	err := service.Delete(context.Background(), id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update_GetByIDErrorAfterUpdate(t *testing.T) {
	mockRepo := new(mocks.CategoryRepositoryMock)
	service := NewCategoryService(mockRepo)

	id := 1
	desc := "Food Updated"
	category := &models.Category{Name: "Food Updated", Description: &desc}

	// Update succeeds
	mockRepo.On("Update", context.Background(), id, category).Return(nil)

	// But GetByID fails (tests lines 46-56 in category_service.go)
	mockRepo.On("GetByID", context.Background(), id).Return(nil, errors.New("database connection lost"))

	result, err := service.Update(context.Background(), id, category)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database connection lost", err.Error())
	mockRepo.AssertExpectations(t)
}

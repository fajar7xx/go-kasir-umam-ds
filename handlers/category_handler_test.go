package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fajar7xx/go-kasir-umam-ds/internal/mocks"
	"fajar7xx/go-kasir-umam-ds/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryHandler_GetAll(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	now := time.Now()
	expectedCategories := []models.Category{
		{ID: 1, Name: "Food", CreatedAt: now},
	}

	mockService.On("GetAll").Return(expectedCategories, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].([]interface{})
	assert.Len(t, data, 1)
	assert.Equal(t, "Food", data[0].(map[string]interface{})["name"])
	mockService.AssertExpectations(t)
}

func TestCategoryHandler_GetAll_Error(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mockService.On("GetAll").Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestCategoryHandler_GetByID(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	now := time.Now()
	expectedCategory := &models.Category{ID: 1, Name: "Food", CreatedAt: now}

	mockService.On("GetByID", 1).Return(expectedCategory, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories/{id}", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])
	assert.Equal(t, "Food", data["name"])
}

func TestCategoryHandler_GetByID_NotFound(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mockService.On("GetByID", 1).Return(nil, errors.New("not found"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories/{id}", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCategoryHandler_Create(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	newCategory := models.Category{Name: "Food", Description: "Tasty"}

	mockService.On("Create", mock.AnythingOfType("*models.Category")).Return(nil)

	body, _ := json.Marshal(newCategory)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCategoryHandler_Create_InvalidBody(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCategoryHandler_Update(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	updatedCategory := &models.Category{Name: "Food Updated"}
	mockService.On("Update", 1, mock.AnythingOfType("*models.Category")).Return(updatedCategory, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /categories/{id}", handler.Update)

	body, _ := json.Marshal(models.Category{Name: "Food Updated"})
	req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCategoryHandler_Update_Error(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mockService.On("Update", 1, mock.AnythingOfType("*models.Category")).Return(nil, errors.New("failed"))

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /categories/{id}", handler.Update)

	body, _ := json.Marshal(models.Category{Name: "Food Updated"})
	req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCategoryHandler_Delete(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mockService.On("Delete", 1).Return(nil)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /categories/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCategoryHandler_Delete_Error(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mockService.On("Delete", 1).Return(errors.New("failed"))

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /categories/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

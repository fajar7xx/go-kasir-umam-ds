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

func TestProductHandler_GetAll(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	now := time.Now()
	expectedProducts := []models.Product{
		{ID: 1, Name: "Nasi Goreng", CreatedAt: now},
	}

	mockService.On("GetAll").Return(expectedProducts, nil)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].([]interface{})
	assert.Len(t, data, 1)
	assert.Equal(t, "Nasi Goreng", data[0].(map[string]interface{})["name"])
	mockService.AssertExpectations(t)
}

func TestProductHandler_GetAll_Error(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	mockService.On("GetAll").Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestProductHandler_GetByID(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	now := time.Now()
	expectedProduct := &models.Product{ID: 1, Name: "Nasi Goreng", CreatedAt: now}

	mockService.On("GetByID", 1).Return(expectedProduct, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products/{id}", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])
	assert.Equal(t, "Nasi Goreng", data["name"])
}

func TestProductHandler_GetByID_NotFound(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	mockService.On("GetByID", 1).Return(nil, errors.New("not found"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products/{id}", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestProductHandler_Create(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	desc := "Tasty"
	newProduct := models.Product{Name: "Nasi Goreng", Description: &desc, Price: 15000}

	mockService.On("Create", mock.AnythingOfType("*models.Product")).Return(nil)

	body, _ := json.Marshal(newProduct)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestProductHandler_Create_InvalidBody(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestProductHandler_Update(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	updatedProduct := &models.Product{Name: "Nasi Goreng Updated"}
	mockService.On("Update", 1, mock.AnythingOfType("*models.Product")).Return(updatedProduct, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /products/{id}", handler.Update)

	body, _ := json.Marshal(models.Product{Name: "Nasi Goreng Updated"})
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestProductHandler_Update_Error(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	mockService.On("Update", 1, mock.AnythingOfType("*models.Product")).Return(nil, errors.New("failed"))

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /products/{id}", handler.Update)

	body, _ := json.Marshal(models.Product{Name: "Nasi Goreng Updated"})
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestProductHandler_Delete(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	mockService.On("Delete", 1).Return(nil)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /products/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestProductHandler_Delete_Error(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	mockService.On("Delete", 1).Return(errors.New("failed"))

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /products/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

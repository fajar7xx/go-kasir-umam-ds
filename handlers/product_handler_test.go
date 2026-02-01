package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fajar7xx/go-kasir-umam-ds/internal/mocks"
	"fajar7xx/go-kasir-umam-ds/models"
	"fajar7xx/go-kasir-umam-ds/utils"
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
	expectedProducts := []models.ProductResponse{
		{ID: 1, Name: "Nasi Goreng", CreatedAt: now},
	}

	mockService.On("GetAll", mock.Anything).Return(expectedProducts, nil)

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

	mockService.On("GetAll", mock.Anything).Return(nil, errors.New("db error"))

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
	expectedProduct := &models.ProductResponse{ID: 1, Name: "Nasi Goreng", CreatedAt: now}

	mockService.On("GetByID", mock.Anything, 1).Return(expectedProduct, nil)

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

	mockService.On("GetByID", mock.Anything, 1).Return(nil, errors.New("not found"))

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
	newProduct := models.Product{Name: "Nasi Goreng", Description: &desc, Price: 15000, Stock: 10, CategoryID: 1}
	expectedResponse := &models.ProductResponse{ID: 1, Name: "Nasi Goreng", Description: &desc, Price: 15000, Stock: 10, CategoryID: 1}

	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Product")).Return(expectedResponse, nil)

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

	updatedProduct := &models.ProductResponse{Name: "Nasi Goreng Updated", Price: 16000, Stock: 15, CategoryID: 1}
	mockService.On("Update", mock.Anything, 1, mock.AnythingOfType("*models.Product")).Return(updatedProduct, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /products/{id}", handler.Update)

	body, _ := json.Marshal(models.Product{Name: "Nasi Goreng Updated", Price: 16000, Stock: 15, CategoryID: 1})
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestProductHandler_Update_Error(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	mockService.On("Update", mock.Anything, 1, mock.AnythingOfType("*models.Product")).Return(nil, errors.New("failed"))

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

	mockService.On("Delete", mock.Anything, 1).Return(nil)

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

	mockService.On("Delete", mock.Anything, 1).Return(errors.New("failed"))

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /products/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// Router dispatch tests
func TestProductHandler_HandleProducts_GetMethod(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	expectedProducts := []models.ProductResponse{}
	mockService.On("GetAll", mock.Anything).Return(expectedProducts, nil)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler.HandleProducts(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandler_HandleProducts_PostMethod(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	desc := "Delicious product"
	product := &models.Product{
		Name:        "Burger",
		Description: &desc,
		Price:       25000,
		Stock:       10,
		CategoryID:  1,
	}
	expectedResponse := &models.ProductResponse{ID: 1, Name: "Burger", Description: &desc, Price: 25000, Stock: 10, CategoryID: 1}
	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Product")).Return(expectedResponse, nil)

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.HandleProducts(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandler_HandleProducts_MethodNotAllowed(t *testing.T) {
	mockService := new(mocks.ProductServiceMock)
	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/products", nil)
	w := httptest.NewRecorder()

	handler.HandleProducts(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	var response utils.ErrorResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "METHOD_NOT_ALLOWED", response.Error.Code)
}

func TestProductHandler_HandleProductByID_AllMethods(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		setupMock      func(*mocks.ProductServiceMock)
		expectedStatus int
	}{
		{
			name:   "GET method",
			method: http.MethodGet,
			setupMock: func(m *mocks.ProductServiceMock) {
				m.On("GetByID", mock.Anything, 1).Return(&models.ProductResponse{ID: 1, Name: "Burger", Price: 25000, Stock: 10, CategoryID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "PUT method",
			method: http.MethodPut,
			setupMock: func(m *mocks.ProductServiceMock) {
				m.On("Update", mock.Anything, 1, mock.AnythingOfType("*models.Product")).Return(&models.ProductResponse{ID: 1, Name: "Updated", Price: 30000, Stock: 15, CategoryID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "PATCH method",
			method: http.MethodPatch,
			setupMock: func(m *mocks.ProductServiceMock) {
				m.On("Update", mock.Anything, 1, mock.AnythingOfType("*models.Product")).Return(&models.ProductResponse{ID: 1, Name: "Updated", Price: 30000, Stock: 15, CategoryID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "DELETE method",
			method: http.MethodDelete,
			setupMock: func(m *mocks.ProductServiceMock) {
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Method not allowed",
			method:         http.MethodPost,
			setupMock:      func(m *mocks.ProductServiceMock) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.ProductServiceMock)
			handler := NewProductHandler(mockService)
			tt.setupMock(mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /products/{id}", handler.HandleProductByID)
			mux.HandleFunc("PUT /products/{id}", handler.HandleProductByID)
			mux.HandleFunc("PATCH /products/{id}", handler.HandleProductByID)
			mux.HandleFunc("DELETE /products/{id}", handler.HandleProductByID)
			mux.HandleFunc("POST /products/{id}", handler.HandleProductByID)

			var req *http.Request
			if tt.method == http.MethodPut || tt.method == http.MethodPatch {
				desc := "Updated product"
				product := models.Product{
					Name:        "Updated",
					Description: &desc,
					Price:       30000,
					Stock:       15,
					CategoryID:  1,
				}
				bodyBytes, _ := json.Marshal(product)
				req = httptest.NewRequest(tt.method, "/products/1", bytes.NewBuffer(bodyBytes))
			} else {
				req = httptest.NewRequest(tt.method, "/products/1", nil)
			}
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

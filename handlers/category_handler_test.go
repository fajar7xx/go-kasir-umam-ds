package handlers

import (
	"bytes"
	"context"
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

func TestCategoryHandler_GetAll(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	now := time.Now()
	expectedCategories := []models.CategoryResponse{
		{ID: 1, Name: "Food", CreatedAt: now},
	}

	mockService.On("GetAll", mock.Anything, mock.Anything).Return(expectedCategories, nil)

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

	mockService.On("GetAll", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

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
	expectedCategory := &models.CategoryResponse{ID: 1, Name: "Food", CreatedAt: now}

	mockService.On("GetByID", mock.Anything, 1).Return(expectedCategory, nil)

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

	mockService.On("GetByID", mock.Anything, 1).Return(nil, errors.New("not found"))

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

	desc := "Tasty"
	newCategory := models.Category{Name: "Food", Description: &desc}
	expectedResponse := &models.CategoryResponse{ID: 1, Name: "Food", Description: &desc}

	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Category")).Return(expectedResponse, nil)

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

	updatedCategory := &models.CategoryResponse{Name: "Food Updated"}
	mockService.On("Update", mock.Anything, 1, mock.AnythingOfType("*models.Category")).Return(updatedCategory, nil)

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

	mockService.On("Update", mock.Anything, 1, mock.AnythingOfType("*models.Category")).Return(nil, errors.New("failed"))

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

	mockService.On("Delete", mock.Anything, 1).Return(nil)

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

	mockService.On("Delete", mock.Anything, 1).Return(errors.New("failed"))

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /categories/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// Router dispatch tests
func TestCategoryHandler_HandleCategories_GetMethod(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	expectedCategories := []models.CategoryResponse{}
	mockService.On("GetAll", mock.Anything, mock.Anything).Return(expectedCategories, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.HandleCategories(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestCategoryHandler_HandleCategories_PostMethod(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	desc := "Food Category"
	category := &models.Category{Name: "Food", Description: &desc}
	expectedResponse := &models.CategoryResponse{ID: 1, Name: "Food", Description: &desc}
	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Category")).Return(expectedResponse, nil)

	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.HandleCategories(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestCategoryHandler_HandleCategories_MethodNotAllowed(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/categories", nil)
	w := httptest.NewRecorder()

	handler.HandleCategories(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	var response utils.ErrorResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "METHOD_NOT_ALLOWED", response.Error.Code)
}

func TestCategoryHandler_HandleCategoryByID_AllMethods(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		setupMock      func(*mocks.CategoryServiceMock)
		expectedStatus int
	}{
		{
			name:   "GET method",
			method: http.MethodGet,
			setupMock: func(m *mocks.CategoryServiceMock) {
				m.On("GetByID", mock.Anything, 1).Return(&models.CategoryResponse{ID: 1, Name: "Food"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "PUT method",
			method: http.MethodPut,
			setupMock: func(m *mocks.CategoryServiceMock) {
				m.On("Update", mock.Anything, 1, mock.AnythingOfType("*models.Category")).Return(&models.CategoryResponse{ID: 1, Name: "Updated"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "PATCH method",
			method: http.MethodPatch,
			setupMock: func(m *mocks.CategoryServiceMock) {
				m.On("Update", mock.Anything, 1, mock.AnythingOfType("*models.Category")).Return(&models.CategoryResponse{ID: 1, Name: "Updated"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "DELETE method",
			method: http.MethodDelete,
			setupMock: func(m *mocks.CategoryServiceMock) {
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Method not allowed",
			method:         http.MethodPost,
			setupMock:      func(m *mocks.CategoryServiceMock) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.CategoryServiceMock)
			handler := NewCategoryHandler(mockService)
			tt.setupMock(mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /categories/{id}", handler.HandleCategoryByID)
			mux.HandleFunc("PUT /categories/{id}", handler.HandleCategoryByID)
			mux.HandleFunc("PATCH /categories/{id}", handler.HandleCategoryByID)
			mux.HandleFunc("DELETE /categories/{id}", handler.HandleCategoryByID)
			mux.HandleFunc("POST /categories/{id}", handler.HandleCategoryByID)

			var req *http.Request
			if tt.method == http.MethodPut || tt.method == http.MethodPatch {
				desc := "Updated"
				category := models.Category{Name: "Updated", Description: &desc}
				bodyBytes, _ := json.Marshal(category)
				req = httptest.NewRequest(tt.method, "/categories/1", bytes.NewBuffer(bodyBytes))
			} else {
				req = httptest.NewRequest(tt.method, "/categories/1", nil)
			}
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

// Validation tests
func TestCategoryHandler_Create_EmptyName(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	category := models.Category{Name: ""}
	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response utils.ErrorResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response.Error.Code, "INVALID")

	mockService.AssertNotCalled(t, "Create")
}

func TestCategoryHandler_Update_EmptyName(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /categories/{id}", handler.HandleCategoryByID)

	category := models.Category{Name: ""}
	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Update")
}

// Context timeout tests
func TestCategoryHandler_GetAll_Timeout(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mockService.On("GetAll", mock.Anything, mock.Anything).Return(nil, context.DeadlineExceeded)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	// When service returns DeadlineExceeded but ctx.Err() is nil, handler returns 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response utils.ErrorResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "INTERNAL_ERROR", response.Error.Code)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_GetByID_Timeout(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mockService.On("GetByID", mock.Anything, 1).Return(nil, context.DeadlineExceeded)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories/{id}", handler.HandleCategoryByID)

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	// Returns 404 when GetByID fails (even with timeout error)
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestCategoryHandler_Create_Timeout(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Category")).Return(nil, context.DeadlineExceeded)

	desc := "Food"
	category := models.Category{Name: "Food", Description: &desc}
	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	// Returns 400 CREATE_FAILED when ctx.Err() is nil but service returns error
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestCategoryHandler_Update_Timeout(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mockService.On("Update", mock.Anything, 1, mock.AnythingOfType("*models.Category")).Return(nil, context.DeadlineExceeded)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /categories/{id}", handler.HandleCategoryByID)

	desc := "Updated"
	category := models.Category{Name: "Updated", Description: &desc}
	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	// Returns 400 when Update fails
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestCategoryHandler_GetByID_InvalidID(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories/{id}", handler.HandleCategoryByID)

	req := httptest.NewRequest(http.MethodGet, "/categories/abc", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response utils.ErrorResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "INVALID_ID", response.Error.Code)

	mockService.AssertNotCalled(t, "GetByID")
}

func TestCategoryHandler_Update_InvalidID(t *testing.T) {
	mockService := new(mocks.CategoryServiceMock)
	handler := NewCategoryHandler(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /categories/{id}", handler.HandleCategoryByID)

	desc := "Food"
	category := models.Category{Name: "Food", Description: &desc}
	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPut, "/categories/xyz", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Update")
}

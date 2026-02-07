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

func TestTransactionHandler_Checkout_Success(t *testing.T) {
	mockService := new(mocks.TransactionServiceMock)
	handler := NewTransactionHandler(mockService)

	now := time.Now()
	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
		{ProductID: 2, Quantity: 1},
	}

	expectedTransaction := &models.Transaction{
		ID:          1,
		TotalAmount: 50000.0,
		Details: []models.TransactionDetail{
			{
				ID:            1,
				TransactionID: 1,
				ProductID:     1,
				ProductName:   "Nasi Goreng",
				Price:         20000.0,
				Quantity:      2,
				SubTotal:      40000.0,
				CreatedAt:     now,
			},
			{
				ID:            2,
				TransactionID: 1,
				ProductID:     2,
				ProductName:   "Es Teh",
				Price:         5000.0,
				Quantity:      2,
				SubTotal:      10000.0,
				CreatedAt:     now,
			},
		},
		CreatedAt: now,
	}

	mockService.On("Checkout", mock.Anything, items, false).Return(expectedTransaction, nil)

	request := models.CheckoutRequest{Items: items}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/transactions/checkout", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Checkout(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])
	assert.Equal(t, 50000.0, data["total_amount"])

	mockService.AssertExpectations(t)
}

func TestTransactionHandler_Checkout_InvalidJSON(t *testing.T) {
	mockService := new(mocks.TransactionServiceMock)
	handler := NewTransactionHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/transactions/checkout", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.Checkout(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response utils.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "INVALID_REQUEST", response.Error.Code)

	mockService.AssertNotCalled(t, "Checkout")
}

func TestTransactionHandler_Checkout_ServiceError(t *testing.T) {
	mockService := new(mocks.TransactionServiceMock)
	handler := NewTransactionHandler(mockService)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
	}

	mockService.On("Checkout", mock.Anything, items, false).Return(nil, errors.New("insufficient stock"))

	request := models.CheckoutRequest{Items: items}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/transactions/checkout", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Checkout(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response utils.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "TRANSACTION_FAILED", response.Error.Code)

	mockService.AssertExpectations(t)
}

func TestTransactionHandler_Checkout_Timeout(t *testing.T) {
	mockService := new(mocks.TransactionServiceMock)
	handler := NewTransactionHandler(mockService)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
	}

	mockService.On("Checkout", mock.Anything, items, false).Return(nil, context.DeadlineExceeded)

	request := models.CheckoutRequest{Items: items}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/transactions/checkout", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Checkout(w, req)

	resp := w.Result()
	// Note: Handler returns 400 when service returns error but ctx.Err() is nil
	// Only returns 504 when actual context timeout occurs (ctx.Err() == DeadlineExceeded)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestTransactionHandler_HandleCheckout_GetMethod(t *testing.T) {
	mockService := new(mocks.TransactionServiceMock)
	handler := NewTransactionHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/transactions/checkout", nil)
	w := httptest.NewRecorder()

	handler.HandleCheckout(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	var response utils.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "METHOD_NOT_ALLOWED", response.Error.Code)
}

func TestTransactionHandler_HandleCheckout_PutMethod(t *testing.T) {
	mockService := new(mocks.TransactionServiceMock)
	handler := NewTransactionHandler(mockService)

	req := httptest.NewRequest(http.MethodPut, "/transactions/checkout", nil)
	w := httptest.NewRecorder()

	handler.HandleCheckout(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestTransactionHandler_HandleCheckout_PostMethod(t *testing.T) {
	mockService := new(mocks.TransactionServiceMock)
	handler := NewTransactionHandler(mockService)

	now := time.Now()
	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 1},
	}

	expectedTransaction := &models.Transaction{
		ID:          1,
		TotalAmount: 20000.0,
		CreatedAt:   now,
		Details:     []models.TransactionDetail{},
	}

	mockService.On("Checkout", mock.Anything, items, false).Return(expectedTransaction, nil)

	request := models.CheckoutRequest{Items: items}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/transactions/checkout", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.HandleCheckout(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

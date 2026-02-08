package handlers

import (
	"encoding/json"
	"errors"
	"fajar7xx/go-kasir-umam-ds/internal/mocks"
	"fajar7xx/go-kasir-umam-ds/models"
	"fajar7xx/go-kasir-umam-ds/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReportHandler_GetReport_Success(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports", nil)
	w := httptest.NewRecorder()

	handler.GetReport(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Report fetched successfully", data["message"])
}

func TestReportHandler_GetReportByPeriod_HariIni_Success(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	expectedReport := &models.ReportResponse{
		TotalRevenue:      350000.0,
		TotalTransactions: 15,
		BestSeller: &models.BestSeller{
			Name:         "Kopi Susu",
			QuantitySold: 25,
		},
	}

	mockService.On("GetTodayReport", mock.Anything).Return(expectedReport, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/hari-ini", nil)
	req.SetPathValue("period", "hari-ini")
	w := httptest.NewRecorder()

	handler.GetReportByPeriod(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, 350000.0, data["total_revenue"])
	assert.Equal(t, float64(15), data["total_transactions"])

	bestSeller := data["best_seller"].(map[string]interface{})
	assert.Equal(t, "Kopi Susu", bestSeller["name"])
	assert.Equal(t, float64(25), bestSeller["quantity_sold"])

	mockService.AssertExpectations(t)
}

func TestReportHandler_GetReportByPeriod_Today_Success(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	expectedReport := &models.ReportResponse{
		TotalRevenue:      180000.0,
		TotalTransactions: 8,
		BestSeller: &models.BestSeller{
			Name:         "Ayam Goreng",
			QuantitySold: 12,
		},
	}

	mockService.On("GetTodayReport", mock.Anything).Return(expectedReport, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/today", nil)
	req.SetPathValue("period", "today")
	w := httptest.NewRecorder()

	handler.GetReportByPeriod(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestReportHandler_GetReportByPeriod_ServiceError(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	mockService.On("GetTodayReport", mock.Anything).
		Return(nil, errors.New("database timeout"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/hari-ini", nil)
	req.SetPathValue("period", "hari-ini")
	w := httptest.NewRecorder()

	handler.GetReportByPeriod(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var response utils.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "INTERNAL_ERROR", response.Error.Code)
	assert.Contains(t, response.Error.Message, "database timeout")

	mockService.AssertExpectations(t)
}

func TestReportHandler_GetReportByPeriod_UnsupportedPeriod(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/weekly", nil)
	req.SetPathValue("period", "weekly")
	w := httptest.NewRecorder()

	handler.GetReportByPeriod(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var response utils.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "NOT_FOUND", response.Error.Code)
	assert.Contains(t, response.Error.Message, "Period not supported")

	mockService.AssertNotCalled(t, "GetTodayReport")
}

func TestReportHandler_HandleReport_MethodNotAllowed(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reports", nil)
	w := httptest.NewRecorder()

	handler.HandleReport(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	var response utils.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "METHOD_NOT_ALLOWED", response.Error.Code)
}

func TestReportHandler_HandleReportByPeriod_MethodNotAllowed(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reports/hari-ini", nil)
	w := httptest.NewRecorder()

	handler.HandleReportByPeriod(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	var response utils.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "METHOD_NOT_ALLOWED", response.Error.Code)
}

func TestReportHandler_GetReportByPeriod_NoBestSeller(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	// Report with no products sold
	expectedReport := &models.ReportResponse{
		TotalRevenue:      0.0,
		TotalTransactions: 0,
		BestSeller:        nil, // No best seller
	}

	mockService.On("GetTodayReport", mock.Anything).Return(expectedReport, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/hari-ini", nil)
	req.SetPathValue("period", "hari-ini")
	w := httptest.NewRecorder()

	handler.GetReportByPeriod(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, 0.0, data["total_revenue"])
	assert.Equal(t, float64(0), data["total_transactions"])
	assert.Nil(t, data["best_seller"]) // Should be nil in JSON

	mockService.AssertExpectations(t)
}

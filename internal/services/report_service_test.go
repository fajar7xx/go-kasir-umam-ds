package services

import (
	"context"
	"errors"
	"fajar7xx/go-kasir-umam-ds/internal/mocks"
	"fajar7xx/go-kasir-umam-ds/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReportService_GetTodayReport_Success(t *testing.T) {
	mockRepo := new(mocks.ReportRepositoryMock)
	service := NewReportService(mockRepo)

	expectedReport := &models.ReportResponse{
		TotalRevenue:      250000.0,
		TotalTransactions: 12,
		BestSeller: &models.BestSeller{
			Name:         "Nasi Goreng Spesial",
			QuantitySold: 20,
		},
	}

	// Mock expects GetDailyReport to be called with any date
	mockRepo.On("GetDailyReport", mock.Anything, mock.Anything).Return(expectedReport, nil)

	result, err := service.GetTodayReport(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 250000.0, result.TotalRevenue)
	assert.Equal(t, 12, result.TotalTransactions)
	assert.NotNil(t, result.BestSeller)
	assert.Equal(t, "Nasi Goreng Spesial", result.BestSeller.Name)
	mockRepo.AssertExpectations(t)
}

func TestReportService_GetReportByDateRange_Success(t *testing.T) {
	mockRepo := new(mocks.ReportRepositoryMock)
	service := NewReportService(mockRepo)

	startDate := "2026-02-01"
	endDate := "2026-02-08"
	expectedStart := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC) // +24 hours

	expectedReport := &models.ReportResponse{
		TotalRevenue:      1500000.0,
		TotalTransactions: 45,
		BestSeller: &models.BestSeller{
			Name:         "Es Teh",
			QuantitySold: 80,
		},
	}

	mockRepo.On("GetReportByDateRange", mock.Anything, expectedStart, expectedEnd).
		Return(expectedReport, nil)

	result, err := service.GetReportByDateRange(context.Background(), startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1500000.0, result.TotalRevenue)
	assert.Equal(t, 45, result.TotalTransactions)
	assert.Equal(t, "Es Teh", result.BestSeller.Name)
	mockRepo.AssertExpectations(t)
}

func TestReportService_GetTodayReport_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.ReportRepositoryMock)
	service := NewReportService(mockRepo)

	expectedError := errors.New("database connection timeout")

	mockRepo.On("GetDailyReport", mock.Anything, mock.Anything).
		Return(nil, expectedError)

	result, err := service.GetTodayReport(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestReportService_GetReportByDateRange_InvalidStartDate(t *testing.T) {
	mockRepo := new(mocks.ReportRepositoryMock)
	service := NewReportService(mockRepo)

	result, err := service.GetReportByDateRange(context.Background(), "invalid-date", "2026-02-08")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid start_date format")
	mockRepo.AssertNotCalled(t, "GetReportByDateRange")
}

func TestReportService_GetReportByDateRange_InvalidEndDate(t *testing.T) {
	mockRepo := new(mocks.ReportRepositoryMock)
	service := NewReportService(mockRepo)

	result, err := service.GetReportByDateRange(context.Background(), "2026-02-01", "not-a-date")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid end_date format")
	mockRepo.AssertNotCalled(t, "GetReportByDateRange")
}

func TestReportService_GetReportByDateRange_EndBeforeStart(t *testing.T) {
	mockRepo := new(mocks.ReportRepositoryMock)
	service := NewReportService(mockRepo)

	result, err := service.GetReportByDateRange(context.Background(), "2026-02-08", "2026-02-01")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "end_date must be after start_date")
	mockRepo.AssertNotCalled(t, "GetReportByDateRange")
}

func TestReportService_GetReportByDateRange_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.ReportRepositoryMock)
	service := NewReportService(mockRepo)

	expectedError := errors.New("query execution failed")

	mockRepo.On("GetReportByDateRange", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, expectedError)

	result, err := service.GetReportByDateRange(context.Background(), "2026-02-01", "2026-02-08")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

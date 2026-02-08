package services

import (
	"context"
	"errors"
	"fajar7xx/go-kasir-umam-ds/internal/repositories"
	"fajar7xx/go-kasir-umam-ds/models"
	"fmt"
	"time"
)

// Custom error types for report validation
var (
	ErrInvalidDateFormat = errors.New("invalid date format")
	ErrInvalidDateRange  = errors.New("end_date must be after start_date")
)

type ReportService interface {
	GetTodayReport(ctx context.Context) (*models.ReportResponse, error)
	GetReportByDateRange(ctx context.Context, startDate, endDate string) (*models.ReportResponse, error)
}

type reportService struct {
	reportRepo repositories.ReportRepository
}

func NewReportService(reportRepo repositories.ReportRepository) ReportService {
	return &reportService{
		reportRepo: reportRepo,
	}
}

func (s *reportService) GetTodayReport(ctx context.Context) (*models.ReportResponse, error) {
	now := time.Now()
	return s.reportRepo.GetDailyReport(ctx, now)
}

func (s *reportService) GetReportByDateRange(ctx context.Context, startDate, endDate string) (*models.ReportResponse, error) {
	layout := "2006-01-02"

	start, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, fmt.Errorf("%w: start_date must be YYYY-MM-DD format", ErrInvalidDateFormat)
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, fmt.Errorf("%w: end_date must be YYYY-MM-DD format", ErrInvalidDateFormat)
	}

	if end.Before(start) {
		return nil, ErrInvalidDateRange
	}

	end = end.Add(24 * time.Hour)

	return s.reportRepo.GetReportByDateRange(ctx, start, end)
}

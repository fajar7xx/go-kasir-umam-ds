package mocks

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/models"
	"time"

	"github.com/stretchr/testify/mock"
)

type ReportRepositoryMock struct {
	mock.Mock
}

func (m *ReportRepositoryMock) GetDailyReport(ctx context.Context, date time.Time) (*models.ReportResponse, error) {
	args := m.Called(ctx, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ReportResponse), args.Error(1)
}

func (m *ReportRepositoryMock) GetReportByDateRange(ctx context.Context, startDate, endDate time.Time) (*models.ReportResponse, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ReportResponse), args.Error(1)
}

package mocks

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/models"

	"github.com/stretchr/testify/mock"
)

type ReportServiceMock struct {
	mock.Mock
}

func (m *ReportServiceMock) GetTodayReport(ctx context.Context) (*models.ReportResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ReportResponse), args.Error(1)
}

func (m *ReportServiceMock) GetReportByDateRange(ctx context.Context, startDate, endDate string) (*models.ReportResponse, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ReportResponse), args.Error(1)
}

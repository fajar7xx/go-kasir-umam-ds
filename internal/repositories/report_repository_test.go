package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestReportRepository_GetReportByDateRange_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewReportRepository(db)

	startDate := time.Date(2026, 2, 8, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC)

	// Mock aggregate query
	aggregateRows := sqlmock.NewRows([]string{"total_revenue", "total_transaksi"}).
		AddRow(150000.0, 5)
	mock.ExpectQuery("SELECT(.+)FROM transactions(.+)").
		WithArgs(startDate, endDate).
		WillReturnRows(aggregateRows)

	// Mock top product query
	topProductRows := sqlmock.NewRows([]string{"name", "qty_terjual"}).
		AddRow("Nasi Goreng", 10)
	mock.ExpectQuery("SELECT(.+)FROM transaction_details(.+)").
		WithArgs(startDate, endDate).
		WillReturnRows(topProductRows)

	result, err := repo.GetReportByDateRange(context.Background(), startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 150000.0, result.TotalRevenue)
	assert.Equal(t, 5, result.TotalTransactions)
	assert.NotNil(t, result.BestSeller)
	assert.Equal(t, "Nasi Goreng", result.BestSeller.Name)
	assert.Equal(t, 10, result.BestSeller.QuantitySold)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_GetReportByDateRange_NoProductsSold(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewReportRepository(db)

	startDate := time.Date(2026, 2, 8, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC)

	// Mock aggregate query - transactions exist but no details
	aggregateRows := sqlmock.NewRows([]string{"total_revenue", "total_transaksi"}).
		AddRow(0.0, 0)
	mock.ExpectQuery("SELECT(.+)FROM transactions(.+)").
		WithArgs(startDate, endDate).
		WillReturnRows(aggregateRows)

	// Mock top product query - no rows returned
	topProductRows := sqlmock.NewRows([]string{"name", "qty_terjual"})
	mock.ExpectQuery("SELECT(.+)FROM transaction_details(.+)").
		WithArgs(startDate, endDate).
		WillReturnRows(topProductRows)

	result, err := repo.GetReportByDateRange(context.Background(), startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0.0, result.TotalRevenue)
	assert.Equal(t, 0, result.TotalTransactions)
	assert.Nil(t, result.BestSeller) // No best seller when no products sold
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_GetReportByDateRange_AggregateQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewReportRepository(db)

	startDate := time.Date(2026, 2, 8, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC)

	// Mock aggregate query error
	mock.ExpectQuery("SELECT(.+)FROM transactions(.+)").
		WithArgs(startDate, endDate).
		WillReturnError(errors.New("database connection failed"))

	result, err := repo.GetReportByDateRange(context.Background(), startDate, endDate)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "aggregate query failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_GetReportByDateRange_TopProductQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewReportRepository(db)

	startDate := time.Date(2026, 2, 8, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC)

	// Mock aggregate query success
	aggregateRows := sqlmock.NewRows([]string{"total_revenue", "total_transaksi"}).
		AddRow(150000.0, 5)
	mock.ExpectQuery("SELECT(.+)FROM transactions(.+)").
		WithArgs(startDate, endDate).
		WillReturnRows(aggregateRows)

	// Mock top product query error
	mock.ExpectQuery("SELECT(.+)FROM transaction_details(.+)").
		WithArgs(startDate, endDate).
		WillReturnError(errors.New("join operation failed"))

	result, err := repo.GetReportByDateRange(context.Background(), startDate, endDate)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "top product query failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_GetDailyReport_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewReportRepository(db)

	date := time.Date(2026, 2, 8, 15, 30, 0, 0, time.UTC)
	expectedStart := time.Date(2026, 2, 8, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC)

	// Mock aggregate query
	aggregateRows := sqlmock.NewRows([]string{"total_revenue", "total_transaksi"}).
		AddRow(200000.0, 8)
	mock.ExpectQuery("SELECT(.+)FROM transactions(.+)").
		WithArgs(expectedStart, expectedEnd).
		WillReturnRows(aggregateRows)

	// Mock top product query
	topProductRows := sqlmock.NewRows([]string{"name", "qty_terjual"}).
		AddRow("Es Teh Manis", 15)
	mock.ExpectQuery("SELECT(.+)FROM transaction_details(.+)").
		WithArgs(expectedStart, expectedEnd).
		WillReturnRows(topProductRows)

	result, err := repo.GetDailyReport(context.Background(), date)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 200000.0, result.TotalRevenue)
	assert.Equal(t, 8, result.TotalTransactions)
	assert.NotNil(t, result.BestSeller)
	assert.Equal(t, "Es Teh Manis", result.BestSeller.Name)
	assert.Equal(t, 15, result.BestSeller.QuantitySold)
	assert.NoError(t, mock.ExpectationsWereMet())
}

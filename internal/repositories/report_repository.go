package repositories

import (
	"context"
	"database/sql"
	"fajar7xx/go-kasir-umam-ds/models"
	"fmt"
	"time"
)

type ReportRepository interface {
	GetDailyReport(ctx context.Context, date time.Time) (*models.ReportResponse, error)
	GetReportByDateRange(ctx context.Context, startDate, endDate time.Time) (*models.ReportResponse, error)
}

type reportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) ReportRepository {
	return &reportRepository{
		db: db,
	}
}

func (r *reportRepository) GetDailyReport(ctx context.Context, date time.Time) (*models.ReportResponse, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)
	return r.GetReportByDateRange(ctx, startOfDay, endOfDay)
}

func (r *reportRepository) GetReportByDateRange(ctx context.Context, startDate, endDate time.Time) (*models.ReportResponse, error) {
	var totalRevenue float64
	var totalTransactions int

	queryAggregate := `
        SELECT
            COALESCE(SUM(total_amount), 0) as total_revenue,
            COUNT(id) as total_transaksi
        FROM transactions
        WHERE created_at >= $1 AND created_at < $2`

	err := r.db.QueryRowContext(ctx, queryAggregate, startDate, endDate).Scan(
		&totalRevenue,
		&totalTransactions,
	)
	if err != nil {
		return nil, fmt.Errorf("aggregate query failed: %w", err)
	}

	// FIX: Initialize bestSeller before scanning
	bestSeller := &models.BestSeller{}
	queryTopProduct := `
        SELECT
            p.name,
            SUM(td.quantity) as qty_terjual
        FROM transaction_details td
        JOIN products p ON td.product_id = p.id
        JOIN transactions t ON td.transaction_id = t.id
        WHERE t.created_at >= $1 AND t.created_at < $2
        GROUP BY p.id, p.name
        ORDER BY qty_terjual DESC
        LIMIT 1`

	err = r.db.QueryRowContext(ctx, queryTopProduct, startDate, endDate).Scan(
		&bestSeller.Name,
		&bestSeller.QuantitySold,
	)
	if err != nil {
		// If no products sold, return nil for bestSeller (not an error)
		if err == sql.ErrNoRows {
			return &models.ReportResponse{
				TotalRevenue:      totalRevenue,
				TotalTransactions: totalTransactions,
				BestSeller:        nil,
			}, nil
		}
		return nil, fmt.Errorf("top product query failed: %w", err)
	}

	return &models.ReportResponse{
		TotalRevenue:      totalRevenue,
		TotalTransactions: totalTransactions,
		BestSeller:        bestSeller,
	}, nil
}

// internal/repositories/transaction_repository_benchmark_test.go
package repositories

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/models"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func BenchmarkCreateTransaction_1Item(b *testing.B) {
	benchmarkCreateTransaction(b, 1)
}

func BenchmarkCreateTransaction_10Items(b *testing.B) {
	benchmarkCreateTransaction(b, 10)
}

func BenchmarkCreateTransaction_50Items(b *testing.B) {
	benchmarkCreateTransaction(b, 50)
}

func BenchmarkCreateTransaction_100Items(b *testing.B) {
	benchmarkCreateTransaction(b, 100)
}

func benchmarkCreateTransaction(b *testing.B, itemCount int) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)

	// Prepare items
	items := make([]models.CheckoutItem, itemCount)
	for i := 0; i < itemCount; i++ {
		items[i] = models.CheckoutItem{
			ProductID: i + 1,
			Quantity:  1,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		setupMockExpectations(mock, items)
		b.StartTimer()

		_, err := repo.CreateTransaction(context.Background(), items)
		if err != nil {
			b.Fatalf("CreateTransaction failed: %v", err)
		}
	}
}

func setupMockExpectations(mock sqlmock.Sqlmock, items []models.CheckoutItem) {
	mock.ExpectBegin()

	// Expect product queries
	for _, item := range items {
		productRows := sqlmock.NewRows([]string{"id", "name", "price", "stock"}).
			AddRow(item.ProductID, fmt.Sprintf("Product %d", item.ProductID), 10000.0, 100)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, stock FROM products WHERE id=$1`)).
			WithArgs(item.ProductID).
			WillReturnRows(productRows)
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2`)).
			WithArgs(item.Quantity, item.ProductID).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	// Expect transaction insert
	now := time.Now()
	totalAmount := float64(len(items)) * 10000.0
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at, updated_at`)).
		WithArgs(totalAmount).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, now, now))

	// Expect batch insert for details
	txIDs := make([]int, len(items))
	productIDs := make([]int, len(items))
	prices := make([]float64, len(items))
	qtys := make([]int, len(items))
	subtotals := make([]float64, len(items))

	for i, item := range items {
		txIDs[i] = 1
		productIDs[i] = item.ProductID
		prices[i] = 10000.0
		qtys[i] = item.Quantity
		subtotals[i] = 10000.0
	}

	detailRows := sqlmock.NewRows([]string{"id", "transaction_id", "product_id", "price", "quantity", "subtotal", "created_at", "updated_at"})
	for i := range items {
		detailRows.AddRow(i+1, 1, productIDs[i], prices[i], qtys[i], subtotals[i], now, now)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal) SELECT * FROM UNNEST( $1::int[], $2::int[], $3::float8[], $4::int[], $5::float8[] ) RETURNING id, transaction_id, product_id, price, quantity, subtotal, created_at, updated_at`)).
		WithArgs(
			pq.Array(txIDs),
			pq.Array(productIDs),
			pq.Array(prices),
			pq.Array(qtys),
			pq.Array(subtotals),
		).
		WillReturnRows(detailRows)

	mock.ExpectCommit()
}

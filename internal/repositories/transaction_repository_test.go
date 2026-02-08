package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fajar7xx/go-kasir-umam-ds/models"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// CreateTransaction Tests (Standard Implementation)
// ============================================================================

func TestTransactionRepository_CreateTransaction_SingleItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
	}

	// Expect Begin
	mock.ExpectBegin()

	// Expect product query
	productRows := sqlmock.NewRows([]string{"id", "name", "price", "stock"}).
		AddRow(1, "Nasi Goreng", 20000.0, 10)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, stock FROM products WHERE id=$1`)).
		WithArgs(1).
		WillReturnRows(productRows)

	// Expect stock update
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2`)).
		WithArgs(2, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Expect transaction insert
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at, updated_at`)).
		WithArgs(40000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, now, now))

	// Expect transaction detail batch insert
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal) SELECT * FROM UNNEST( $1::int[], $2::int[], $3::float8[], $4::int[], $5::float8[] ) RETURNING id, transaction_id, product_id, price, quantity, subtotal, created_at, updated_at`)).
		WithArgs(
			pq.Array([]int{1}),
			pq.Array([]int{1}),
			pq.Array([]float64{20000.0}),
			pq.Array([]int{2}),
			pq.Array([]float64{40000.0}),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "transaction_id", "product_id", "price", "quantity", "subtotal", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 20000.0, 2, 40000.0, now, now))

	// Expect Commit
	mock.ExpectCommit()

	result, err := repo.CreateTransaction(context.Background(), items)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 40000.0, result.TotalAmount)
	assert.Len(t, result.Details, 1)
	assert.Equal(t, "Nasi Goreng", result.Details[0].ProductName)

	// Verify timestamps are populated
	assert.NotZero(t, result.CreatedAt)
	assert.NotNil(t, result.UpdatedAt)
	assert.NotZero(t, *result.UpdatedAt)
	assert.NotZero(t, result.Details[0].ID)
	assert.NotZero(t, result.Details[0].CreatedAt)
	assert.NotNil(t, result.Details[0].UpdatedAt)
	assert.NotZero(t, *result.Details[0].UpdatedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_CreateTransaction_ProductNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)

	items := []models.CheckoutItem{
		{ProductID: 999, Quantity: 1},
	}

	// Expect Begin
	mock.ExpectBegin()

	// Expect product query - returns no rows
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, stock FROM products WHERE id=$1`)).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	// Expect Rollback
	mock.ExpectRollback()

	result, err := repo.CreateTransaction(context.Background(), items)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "product not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_CreateTransaction_BeginError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 1},
	}

	// Expect Begin to fail
	mock.ExpectBegin().WillReturnError(errors.New("connection error"))

	result, err := repo.CreateTransaction(context.Background(), items)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error begin transaction")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_CreateTransaction_BatchInsertDetails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
		{ProductID: 2, Quantity: 3},
		{ProductID: 3, Quantity: 1},
	}

	// Expect Begin
	mock.ExpectBegin()

	// Expect individual product queries (kept as-is)
	for _, item := range items {
		productRows := sqlmock.NewRows([]string{"id", "name", "price", "stock"}).
			AddRow(item.ProductID, fmt.Sprintf("Product %d", item.ProductID), float64(item.ProductID*5000), 100)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, stock FROM products WHERE id=$1`)).
			WithArgs(item.ProductID).
			WillReturnRows(productRows)
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2`)).
			WithArgs(item.Quantity, item.ProductID).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	// Expect transaction insert
	now := time.Now()
	totalAmount := float64(2*5000 + 3*10000 + 1*15000) // 10000 + 30000 + 15000 = 55000
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at, updated_at`)).
		WithArgs(totalAmount).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, now, now))

	// Expect BATCH insert for transaction details using UNNEST
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal) SELECT * FROM UNNEST( $1::int[], $2::int[], $3::float8[], $4::int[], $5::float8[] ) RETURNING id, transaction_id, product_id, price, quantity, subtotal, created_at, updated_at`)).
		WithArgs(
			pq.Array([]int{1, 1, 1}),                       // transaction_id for all 3 items
			pq.Array([]int{1, 2, 3}),                       // product_ids
			pq.Array([]float64{5000.0, 10000.0, 15000.0}),  // prices
			pq.Array([]int{2, 3, 1}),                       // quantities
			pq.Array([]float64{10000.0, 30000.0, 15000.0}), // subtotals
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "transaction_id", "product_id", "price", "quantity", "subtotal", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 5000.0, 2, 10000.0, now, now).
			AddRow(2, 1, 2, 10000.0, 3, 30000.0, now, now).
			AddRow(3, 1, 3, 15000.0, 1, 15000.0, now, now))

	// Expect Commit
	mock.ExpectCommit()

	result, err := repo.CreateTransaction(context.Background(), items)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 55000.0, result.TotalAmount)
	assert.Len(t, result.Details, 3)

	// Verify all details have database values
	for i, detail := range result.Details {
		assert.NotZero(t, detail.ID, "detail %d should have ID", i)
		assert.NotZero(t, detail.CreatedAt, "detail %d should have CreatedAt", i)
		assert.NotNil(t, detail.UpdatedAt, "detail %d should have UpdatedAt", i)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================================
// CreateTransactionOptimal Tests (Batch Implementation)
// ============================================================================

func TestTransactionRepository_CreateTransactionOptimal_SingleItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
	}

	// Expect BeginTx
	mock.ExpectBegin()

	// Expect batch SELECT with FOR UPDATE
	productRows := sqlmock.NewRows([]string{"id", "name", "price", "stock"}).
		AddRow(1, "Nasi Goreng", 20000.0, 10)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, stock FROM products WHERE id IN ($1) FOR UPDATE`)).
		WithArgs(1).
		WillReturnRows(productRows)

	// Expect batch UPDATE with pq.Array
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE products AS p SET stock = p.stock - u.qty, updated_at = NOW() FROM ( SELECT UNNEST($1::int[]) AS id, UNNEST($2::int[]) AS qty ) AS u WHERE p.id = u.id`)).
		WithArgs(pq.Array([]int{1}), pq.Array([]int{2})).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Expect transaction insert
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at, updated_at`)).
		WithArgs(40000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, now, now))

	// Expect batch INSERT transaction details with UNNEST
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal) SELECT * FROM UNNEST( $1::int[], $2::int[], $3::float8[], $4::int[], $5::float8[] ) RETURNING id, transaction_id, product_id, price, quantity, subtotal, created_at, updated_at`)).
		WithArgs(
			pq.Array([]int{1}),
			pq.Array([]int{1}),
			pq.Array([]float64{20000.0}),
			pq.Array([]int{2}),
			pq.Array([]float64{40000.0}),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "transaction_id", "product_id", "price", "quantity", "subtotal", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 20000.0, 2, 40000.0, now, now))

	// Expect Commit
	mock.ExpectCommit()

	result, err := repo.CreateTransactionOptimal(context.Background(), items)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 40000.0, result.TotalAmount)
	assert.Len(t, result.Details, 1)

	// Verify timestamps are populated
	assert.NotZero(t, result.CreatedAt)
	assert.NotNil(t, result.UpdatedAt)
	assert.NotZero(t, *result.UpdatedAt)
	assert.NotZero(t, result.Details[0].ID)
	assert.NotZero(t, result.Details[0].CreatedAt)
	assert.NotNil(t, result.Details[0].UpdatedAt)
	assert.NotZero(t, *result.Details[0].UpdatedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_CreateTransactionOptimal_ProductNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)

	items := []models.CheckoutItem{
		{ProductID: 999, Quantity: 1},
	}

	// Expect BeginTx
	mock.ExpectBegin()

	// Expect batch SELECT - returns no rows
	productRows := sqlmock.NewRows([]string{"id", "name", "price", "stock"})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, stock FROM products WHERE id IN ($1) FOR UPDATE`)).
		WithArgs(999).
		WillReturnRows(productRows)

	// Expect Rollback
	mock.ExpectRollback()

	result, err := repo.CreateTransactionOptimal(context.Background(), items)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "one or more products not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_CreateTransactionOptimal_InsufficientStock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 100}, // Request more than available
	}

	// Expect BeginTx
	mock.ExpectBegin()

	// Expect batch SELECT - returns product with insufficient stock
	productRows := sqlmock.NewRows([]string{"id", "name", "price", "stock"}).
		AddRow(1, "Nasi Goreng", 20000.0, 5) // Only 5 in stock
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, stock FROM products WHERE id IN ($1) FOR UPDATE`)).
		WithArgs(1).
		WillReturnRows(productRows)

	// Expect Rollback
	mock.ExpectRollback()

	result, err := repo.CreateTransactionOptimal(context.Background(), items)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "insufficient stock")
	assert.NoError(t, mock.ExpectationsWereMet())
}

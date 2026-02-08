# Optimize Transaction Details Insertion - Performance Fix

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace loop-based individual INSERTs with a single batch INSERT using UNNEST for transaction_details in CreateTransaction method

**Architecture:** Change from N individual INSERT queries to 1 batch INSERT query using PostgreSQL's UNNEST function, reducing database round-trips

**Tech Stack:** Go 1.22+, PostgreSQL, database/sql, lib/pq, sqlmock, testify

---

## Background

**Current Problem:**
- `CreateTransaction` inserts transaction details one-by-one in a loop
- For 10 items = 10 separate INSERT queries
- Each query is a round-trip to the database
- Performance degrades with more items

**Current Code (lines 122-150):**
```go
for i := range details {
    details[i].TransactionID = transactionID
    // Individual INSERT for each detail
    err = tx.QueryRowContext(ctx, insertTransactionDetailQuery, ...).Scan(...)
}
```

**Performance Impact:**
- 1 item: 1 INSERT query
- 10 items: 10 INSERT queries  
- 100 items: 100 INSERT queries

**Solution:**
Use PostgreSQL's UNNEST to insert all details in a single query, like `CreateTransactionOptimal` already does.

**Performance Improvement:**
- 1 item: 1 INSERT query
- 10 items: 1 INSERT query ✅
- 100 items: 1 INSERT query ✅

---

## Task 1: Add Helper Function for Batch Insert

**Files:**
- Modify: `internal/repositories/transaction_repository.go:122-150`
- Test: `internal/repositories/transaction_repository_test.go`

**Step 1: Write failing test for batch insert helper**

```go
// Add to internal/repositories/transaction_repository_test.go

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
	totalAmount := float64(1*2*5000 + 2*3*10000 + 3*1*15000) // 10000 + 60000 + 15000 = 85000
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at, updated_at`)).
		WithArgs(85000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, now, now))

	// Expect BATCH insert for transaction details using UNNEST
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal) SELECT * FROM UNNEST( $1::int[], $2::int[], $3::float8[], $4::int[], $5::float8[] ) RETURNING id, transaction_id, product_id, price, quantity, subtotal, created_at, updated_at`)).
		WithArgs(
			pq.Array([]int{1, 1, 1}),                    // transaction_id for all 3 items
			pq.Array([]int{1, 2, 3}),                    // product_ids
			pq.Array([]float64{5000.0, 10000.0, 15000.0}), // prices
			pq.Array([]int{2, 3, 1}),                    // quantities
			pq.Array([]float64{10000.0, 60000.0, 15000.0}), // subtotals
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "transaction_id", "product_id", "price", "quantity", "subtotal", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 5000.0, 2, 10000.0, now, now).
			AddRow(2, 1, 2, 10000.0, 3, 60000.0, now, now).
			AddRow(3, 1, 3, 15000.0, 1, 15000.0, now, now))

	// Expect Commit
	mock.ExpectCommit()

	result, err := repo.CreateTransaction(context.Background(), items)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 85000.0, result.TotalAmount)
	assert.Len(t, result.Details, 3)

	// Verify all details have database values
	for i, detail := range result.Details {
		assert.NotZero(t, detail.ID, "detail %d should have ID", i)
		assert.NotZero(t, detail.CreatedAt, "detail %d should have CreatedAt", i)
		assert.NotNil(t, detail.UpdatedAt, "detail %d should have UpdatedAt", i)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/fajarsiagian/Learning/golang/go-kasir-umam-ds
go test ./internal/repositories -v -run TestTransactionRepository_CreateTransaction_BatchInsertDetails
```

Expected: FAIL - test expects batch INSERT but code still uses loop

**Step 3: Implement batch insert for transaction details**

Replace the loop-based insert (lines 122-150) with batch insert:

```go
// Modify internal/repositories/transaction_repository.go:122-150
// Replace this section:
//   insertTransactionDetailQuery := `INSERT INTO transaction_details ...`
//   for i := range details { ... }
//
// With this:

	// insert transaction_details using BATCH INSERT
	if len(details) > 0 {
		// Prepare arrays for batch insert
		txIDs := make([]int, len(details))
		productIDs := make([]int, len(details))
		prices := make([]float64, len(details))
		quantities := make([]int, len(details))
		subtotals := make([]float64, len(details))

		for i, detail := range details {
			txIDs[i] = transactionID
			productIDs[i] = detail.ProductID
			prices[i] = detail.Price
			quantities[i] = detail.Quantity
			subtotals[i] = detail.SubTotal
		}

		// Single batch INSERT using UNNEST
		insertDetailsQuery := `INSERT INTO transaction_details
								(transaction_id, product_id, price, quantity, subtotal)
								SELECT * FROM UNNEST(
									$1::int[],
									$2::int[],
									$3::float8[],
									$4::int[],
									$5::float8[]
								)
								RETURNING id, transaction_id, product_id, price, quantity, subtotal, created_at, updated_at`

		rows, err := tx.QueryContext(ctx, insertDetailsQuery,
			pq.Array(txIDs),
			pq.Array(productIDs),
			pq.Array(prices),
			pq.Array(quantities),
			pq.Array(subtotals),
		)
		if err != nil {
			return nil, fmt.Errorf("error to insert transaction details: %w", err)
		}
		defer rows.Close()

		// Create map for O(1) product name lookup
		productNameMap := make(map[int]string, len(details))
		for _, detail := range details {
			productNameMap[detail.ProductID] = detail.ProductName
		}

		// Reconstruct details with database values
		finalDetails := make([]models.TransactionDetail, 0, len(details))
		for rows.Next() {
			var d models.TransactionDetail
			err := rows.Scan(
				&d.ID,
				&d.TransactionID,
				&d.ProductID,
				&d.Price,
				&d.Quantity,
				&d.SubTotal,
				&d.CreatedAt,
				&d.UpdatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("error scanning transaction detail: %w", err)
			}

			// Restore product name from map
			d.ProductName = productNameMap[d.ProductID]
			finalDetails = append(finalDetails, d)
		}

		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating transaction details: %w", err)
		}

		// Replace details with finalDetails containing database values
		details = finalDetails
	}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/repositories -v -run TestTransactionRepository_CreateTransaction_BatchInsertDetails
```

Expected: PASS

**Step 5: Run all existing tests to ensure no regression**

```bash
go test ./internal/repositories -v -run TestTransactionRepository_CreateTransaction
```

Expected: ALL PASS (including the original single-item test)

**Step 6: Commit**

```bash
git add internal/repositories/transaction_repository.go internal/repositories/transaction_repository_test.go
git commit -m "perf: optimize CreateTransaction to use batch INSERT for transaction details

- Replace loop-based individual INSERTs with single UNNEST batch INSERT
- Reduces DB queries from N to 1 for transaction_details
- Performance: 10 items now takes 1 query instead of 10
- Maintains all database-generated values (ID, timestamps)
- All existing tests pass"
```

---

## Task 2: Performance Comparison Test (Optional)

**Files:**
- Create: `internal/repositories/transaction_repository_benchmark_test.go`

**Step 1: Create benchmark test**

```go
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
```

**Step 2: Run benchmark**

```bash
go test ./internal/repositories -bench=BenchmarkCreateTransaction -benchmem
```

Expected output shows performance improvement:
```
BenchmarkCreateTransaction_1Item-8       5000    250000 ns/op    15000 B/op    120 allocs/op
BenchmarkCreateTransaction_10Items-8     2000    800000 ns/op    85000 B/op    650 allocs/op
BenchmarkCreateTransaction_50Items-8      500   3500000 ns/op   410000 B/op   3100 allocs/op
BenchmarkCreateTransaction_100Items-8     200   7200000 ns/op   820000 B/op   6200 allocs/op
```

**Step 3: Commit benchmark**

```bash
git add internal/repositories/transaction_repository_benchmark_test.go
git commit -m "test: add benchmark tests for CreateTransaction performance"
```

---

## Task 3: Update Documentation

**Files:**
- Create: `docs/architecture/transaction-performance-optimization.md`
- Modify: `docs/README.md`

**Step 1: Create performance documentation**

```markdown
<!-- docs/architecture/transaction-performance-optimization.md -->
# Transaction Performance Optimization

## Overview

The `CreateTransaction` method has been optimized to use batch INSERT operations for transaction details, significantly reducing database round-trips.

## Problem

**Before Optimization:**

```go
// Loop-based: N separate INSERT queries
for i := range details {
    tx.QueryRowContext(ctx, 
        `INSERT INTO transaction_details (...) VALUES ($1, $2, $3, $4, $5)`,
        transactionID, productID, price, qty, subtotal)
}
```

**Performance:**
- 1 item: 1 INSERT query
- 10 items: 10 INSERT queries
- 100 items: 100 INSERT queries

Each INSERT is a separate database round-trip (network latency + query execution).

## Solution

**After Optimization:**

```go
// Batch INSERT: 1 query for all items
tx.QueryContext(ctx,
    `INSERT INTO transaction_details (...)
     SELECT * FROM UNNEST($1::int[], $2::int[], $3::float8[], $4::int[], $5::float8[])`,
    pq.Array(txIDs), pq.Array(productIDs), pq.Array(prices), 
    pq.Array(quantities), pq.Array(subtotals))
```

**Performance:**
- 1 item: 1 INSERT query
- 10 items: 1 INSERT query ✅
- 100 items: 1 INSERT query ✅

## Implementation Details

### PostgreSQL UNNEST Function

UNNEST takes multiple arrays and combines them row-by-row:

```sql
SELECT * FROM UNNEST(
    ARRAY[1, 1, 1],              -- transaction_ids
    ARRAY[101, 102, 103],        -- product_ids
    ARRAY[5000, 10000, 15000],   -- prices
    ARRAY[2, 3, 1],              -- quantities
    ARRAY[10000, 30000, 15000]   -- subtotals
)
-- Produces 3 rows:
-- (1, 101, 5000, 2, 10000)
-- (1, 102, 10000, 3, 30000)
-- (1, 103, 15000, 1, 15000)
```

### Data Preservation

All database-generated values are preserved using `RETURNING`:

```sql
... RETURNING id, transaction_id, product_id, price, quantity, 
              subtotal, created_at, updated_at
```

The implementation:
1. Scans all returned rows
2. Maps product names back to details
3. Returns complete `TransactionDetail` structs with DB values

### Code Structure

```go
// 1. Prepare arrays
txIDs := make([]int, len(details))
productIDs := make([]int, len(details))
// ... etc

// 2. Batch insert
rows, err := tx.QueryContext(ctx, insertQuery, 
    pq.Array(txIDs), pq.Array(productIDs), ...)

// 3. Scan results
for rows.Next() {
    var d models.TransactionDetail
    rows.Scan(&d.ID, &d.TransactionID, ...)
    finalDetails = append(finalDetails, d)
}

// 4. Restore product names (not in DB)
d.ProductName = productNameMap[d.ProductID]
```

## Performance Impact

### Query Count Reduction

| Items | Queries Before | Queries After | Improvement |
|-------|---------------|---------------|-------------|
| 1     | 1             | 1             | 0%          |
| 10    | 10            | 1             | **90%**     |
| 50    | 50            | 1             | **98%**     |
| 100   | 100           | 1             | **99%**     |

### Real-World Impact

Assuming 5ms per query (network + execution):

| Items | Time Before | Time After | Saved   |
|-------|------------|------------|---------|
| 10    | 50ms       | 5ms        | 45ms    |
| 50    | 250ms      | 5ms        | 245ms   |
| 100   | 500ms      | 5ms        | 495ms   |

**For a transaction with 100 items, we go from 500ms to 5ms - a 100x improvement!**

## When to Use

This optimization is most beneficial when:
- ✅ Transactions have multiple items (5+)
- ✅ Database is remote (network latency matters)
- ✅ High transaction volume

Less impactful when:
- ⚠️ Transactions typically have 1-2 items
- ⚠️ Database is local (localhost)
- ⚠️ Low transaction volume

## Backward Compatibility

- ✅ All existing tests pass
- ✅ API contract unchanged
- ✅ All database values preserved
- ✅ No breaking changes

## Alternative: CreateTransactionOptimal

For even better performance, consider `CreateTransactionOptimal` which also batches:
- Product SELECT queries (single `IN` query)
- Stock UPDATE queries (single `UNNEST` update)
- Row locking with `FOR UPDATE`

See [Transaction Repository Design](transaction-repository-design.md) for details.

## Testing

Run benchmarks to measure performance:

```bash
go test ./internal/repositories -bench=BenchmarkCreateTransaction -benchmem
```

## References

- PostgreSQL UNNEST: https://www.postgresql.org/docs/current/functions-array.html
- lib/pq Arrays: https://pkg.go.dev/github.com/lib/pq#Array
```

**Step 2: Update docs/README.md**

```markdown
<!-- Add to docs/README.md under Architecture section -->

## Architecture Documentation

- [Overview](architecture/overview.md) - High-level system architecture
- [Clean Architecture](architecture/clean-architecture.md) - Clean architecture patterns
- [Database Schema](architecture/database-schema.md) - Database design
- **[Transaction Performance](architecture/transaction-performance-optimization.md)** - Batch INSERT optimization ✨ NEW
```

**Step 3: Commit documentation**

```bash
git add docs/architecture/transaction-performance-optimization.md docs/README.md
git commit -m "docs: add transaction performance optimization documentation"
```

---

## Completion Checklist

- [ ] Task 1: Batch INSERT implementation
- [ ] Task 2: Benchmark tests (optional)
- [ ] Task 3: Documentation

## Success Criteria

✅ Transaction details inserted in 1 query instead of N  
✅ All existing tests pass  
✅ Database-generated values preserved  
✅ No breaking changes  
✅ Performance improvement documented

## Performance Gains

**Before:**
- 10 items = 10 INSERT queries (~50ms)

**After:**
- 10 items = 1 INSERT query (~5ms)

**90% reduction in queries!**

---

**End of Implementation Plan**

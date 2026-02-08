# Transaction Repository Refactoring Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Consolidate duplicate transaction creation methods into a single, configurable implementation with proper service layer integration

**Architecture:** Merge `CreateTransaction` and `CreateTransactionOptimal` into one method that accepts options for batch operations and row locking. Update service layer to properly utilize the `useLock` parameter.

**Tech Stack:** Go 1.22+, PostgreSQL, database/sql, lib/pq, sqlmock, testify

---

## Background

The current implementation has:
- Two separate methods: `CreateTransaction` (loop-based) and `CreateTransactionOptimal` (batch-based)
- Service layer with unused `useLock` parameter
- Duplication of transaction logic across both methods
- Both methods already correctly use RETURNING to capture DB-generated values

**Improvement Goals:**
1. Single source of truth for transaction creation
2. Configurable batch vs. loop operations
3. Proper service layer integration
4. Maintain 100% test coverage

---

## Task 1: Add Transaction Options Type

**Files:**
- Create: `internal/repositories/transaction_options.go`
- Test: `internal/repositories/transaction_options_test.go`

**Step 1: Write failing test for TransactionOptions**

```go
// internal/repositories/transaction_options_test.go
package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionOptions_Default(t *testing.T) {
	opts := DefaultTransactionOptions()
	
	assert.False(t, opts.UseBatchOperations)
	assert.False(t, opts.UseRowLocking)
}

func TestTransactionOptions_WithBatchOperations(t *testing.T) {
	opts := DefaultTransactionOptions()
	opts.UseBatchOperations = true
	
	assert.True(t, opts.UseBatchOperations)
	assert.False(t, opts.UseRowLocking)
}

func TestTransactionOptions_WithRowLocking(t *testing.T) {
	opts := DefaultTransactionOptions()
	opts.UseRowLocking = true
	
	assert.True(t, opts.UseRowLocking)
	assert.False(t, opts.UseBatchOperations)
}

func TestTransactionOptions_BothEnabled(t *testing.T) {
	opts := TransactionOptions{
		UseBatchOperations: true,
		UseRowLocking:      true,
	}
	
	assert.True(t, opts.UseBatchOperations)
	assert.True(t, opts.UseRowLocking)
}
```

**Step 2: Run test to verify it fails**

```bash
cd /home/fajarsiagian/Learning/golang/go-kasir-umam-ds
go test ./internal/repositories -v -run TestTransactionOptions
```

Expected: FAIL - undefined: DefaultTransactionOptions

**Step 3: Implement TransactionOptions**

```go
// internal/repositories/transaction_options.go
package repositories

// TransactionOptions configures how CreateTransaction executes
type TransactionOptions struct {
	// UseBatchOperations enables batch INSERT/UPDATE with UNNEST
	// When false, uses traditional loop-based operations
	UseBatchOperations bool
	
	// UseRowLocking enables SELECT ... FOR UPDATE on products
	// Provides better concurrency control for high-traffic scenarios
	UseRowLocking bool
}

// DefaultTransactionOptions returns conservative defaults
func DefaultTransactionOptions() TransactionOptions {
	return TransactionOptions{
		UseBatchOperations: false,
		UseRowLocking:      false,
	}
}

// OptimalTransactionOptions returns performance-optimized settings
func OptimalTransactionOptions() TransactionOptions {
	return TransactionOptions{
		UseBatchOperations: true,
		UseRowLocking:      true,
	}
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/repositories -v -run TestTransactionOptions
```

Expected: PASS (4/4 tests)

**Step 5: Commit**

```bash
git add internal/repositories/transaction_options.go internal/repositories/transaction_options_test.go
git commit -m "feat: add TransactionOptions for configurable transaction behavior"
```

---

## Task 2: Update Repository Interface

**Files:**
- Modify: `internal/repositories/transaction_repository.go:9-12`
- Test: `internal/repositories/transaction_repository_test.go`

**Step 1: Write failing test for new interface**

```go
// Add to internal/repositories/transaction_repository_test.go
func TestTransactionRepository_CreateTransactionWithOptions_Interface(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)
	
	// Verify interface includes new method
	var _ TransactionRepository = repo
	
	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 1},
	}
	opts := DefaultTransactionOptions()
	
	// This will fail because method doesn't exist yet
	_, err = repo.CreateTransactionWithOptions(context.Background(), items, opts)
	assert.Error(t, err) // Expected to fail - method not implemented
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/repositories -v -run TestTransactionRepository_CreateTransactionWithOptions_Interface
```

Expected: Compilation error - repo.CreateTransactionWithOptions undefined

**Step 3: Update repository interface**

```go
// Modify internal/repositories/transaction_repository.go:9-12
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, items []models.CheckoutItem) (*models.Transaction, error)
	CreateTransactionOptimal(ctx context.Context, items []models.CheckoutItem) (*models.Transaction, error)
	// New unified method
	CreateTransactionWithOptions(ctx context.Context, items []models.CheckoutItem, opts TransactionOptions) (*models.Transaction, error)
}
```

**Step 4: Add stub implementation**

```go
// Add to internal/repositories/transaction_repository.go (after CreateTransactionOptimal)
func (repo *transactionRepository) CreateTransactionWithOptions(ctx context.Context, items []models.CheckoutItem, opts TransactionOptions) (*models.Transaction, error) {
	// Temporary implementation - delegate to existing methods
	if opts.UseBatchOperations {
		return repo.CreateTransactionOptimal(ctx, items)
	}
	return repo.CreateTransaction(ctx, items)
}
```

**Step 5: Update test to expect success**

```go
// Update test in internal/repositories/transaction_repository_test.go
func TestTransactionRepository_CreateTransactionWithOptions_Interface(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)
	
	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 1},
	}
	opts := DefaultTransactionOptions()
	
	// Setup mock expectations for CreateTransaction (default path)
	mock.ExpectBegin()
	productRows := sqlmock.NewRows([]string{"id", "name", "price", "stock"}).
		AddRow(1, "Test Product", 10000.0, 10)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, stock FROM products WHERE id=$1`)).
		WithArgs(1).
		WillReturnRows(productRows)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2`)).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at, updated_at`)).
		WithArgs(10000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, now, now))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`)).
		WithArgs(1, 1, 10000.0, 1, 10000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, now, now))
	mock.ExpectCommit()
	
	result, err := repo.CreateTransactionWithOptions(context.Background(), items, opts)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
```

**Step 6: Run test to verify it passes**

```bash
go test ./internal/repositories -v -run TestTransactionRepository_CreateTransactionWithOptions_Interface
```

Expected: PASS

**Step 7: Commit**

```bash
git add internal/repositories/transaction_repository.go internal/repositories/transaction_repository_test.go
git commit -m "feat: add CreateTransactionWithOptions method to repository interface"
```

---

## Task 3: Implement Unified Transaction Creation Logic

**Files:**
- Modify: `internal/repositories/transaction_repository.go:270-end`
- Test: `internal/repositories/transaction_repository_test.go`

**Step 1: Write comprehensive tests for unified method**

```go
// Add to internal/repositories/transaction_repository_test.go

func TestTransactionRepository_CreateTransactionWithOptions_DefaultOptions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)
	items := []models.CheckoutItem{{ProductID: 1, Quantity: 2}}
	opts := DefaultTransactionOptions()

	// Expect loop-based operations (no batch, no locking)
	mock.ExpectBegin()
	productRows := sqlmock.NewRows([]string{"id", "name", "price", "stock"}).
		AddRow(1, "Product", 5000.0, 10)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, stock FROM products WHERE id=$1`)).
		WithArgs(1).WillReturnRows(productRows)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2`)).
		WithArgs(2, 1).WillReturnResult(sqlmock.NewResult(0, 1))
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at, updated_at`)).
		WithArgs(10000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, now, now))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`)).
		WithArgs(1, 1, 5000.0, 2, 10000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, now, now))
	mock.ExpectCommit()

	result, err := repo.CreateTransactionWithOptions(context.Background(), items, opts)
	assert.NoError(t, err)
	assert.Equal(t, 10000.0, result.TotalAmount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_CreateTransactionWithOptions_BatchOperations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)
	items := []models.CheckoutItem{{ProductID: 1, Quantity: 2}}
	opts := TransactionOptions{UseBatchOperations: true, UseRowLocking: false}

	// Expect batch operations (with UNNEST)
	mock.ExpectBegin()
	productRows := sqlmock.NewRows([]string{"id", "name", "price", "stock"}).
		AddRow(1, "Product", 5000.0, 10)
	// Note: No FOR UPDATE when UseRowLocking is false
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, stock FROM products WHERE id IN ($1)`)).
		WithArgs(1).WillReturnRows(productRows)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE products AS p SET stock = p.stock - u.qty, updated_at = NOW() FROM ( SELECT UNNEST($1::int[]) AS id, UNNEST($2::int[]) AS qty ) AS u WHERE p.id = u.id`)).
		WithArgs(pq.Array([]int{1}), pq.Array([]int{2})).
		WillReturnResult(sqlmock.NewResult(0, 1))
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at, updated_at`)).
		WithArgs(10000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, now, now))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal) SELECT * FROM UNNEST( $1::int[], $2::int[], $3::float8[], $4::int[], $5::float8[] ) RETURNING id, transaction_id, product_id, price, quantity, subtotal, created_at, updated_at`)).
		WithArgs(pq.Array([]int{1}), pq.Array([]int{1}), pq.Array([]float64{5000.0}), pq.Array([]int{2}), pq.Array([]float64{10000.0})).
		WillReturnRows(sqlmock.NewRows([]string{"id", "transaction_id", "product_id", "price", "quantity", "subtotal", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 5000.0, 2, 10000.0, now, now))
	mock.ExpectCommit()

	result, err := repo.CreateTransactionWithOptions(context.Background(), items, opts)
	assert.NoError(t, err)
	assert.Equal(t, 10000.0, result.TotalAmount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_CreateTransactionWithOptions_WithRowLocking(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)
	items := []models.CheckoutItem{{ProductID: 1, Quantity: 2}}
	opts := TransactionOptions{UseBatchOperations: true, UseRowLocking: true}

	// Expect batch operations WITH FOR UPDATE
	mock.ExpectBegin()
	productRows := sqlmock.NewRows([]string{"id", "name", "price", "stock"}).
		AddRow(1, "Product", 5000.0, 10)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, stock FROM products WHERE id IN ($1) FOR UPDATE`)).
		WithArgs(1).WillReturnRows(productRows)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE products AS p SET stock = p.stock - u.qty, updated_at = NOW() FROM ( SELECT UNNEST($1::int[]) AS id, UNNEST($2::int[]) AS qty ) AS u WHERE p.id = u.id`)).
		WithArgs(pq.Array([]int{1}), pq.Array([]int{2})).
		WillReturnResult(sqlmock.NewResult(0, 1))
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at, updated_at`)).
		WithArgs(10000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, now, now))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal) SELECT * FROM UNNEST( $1::int[], $2::int[], $3::float8[], $4::int[], $5::float8[] ) RETURNING id, transaction_id, product_id, price, quantity, subtotal, created_at, updated_at`)).
		WithArgs(pq.Array([]int{1}), pq.Array([]int{1}), pq.Array([]float64{5000.0}), pq.Array([]int{2}), pq.Array([]float64{10000.0})).
		WillReturnRows(sqlmock.NewRows([]string{"id", "transaction_id", "product_id", "price", "quantity", "subtotal", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 5000.0, 2, 10000.0, now, now))
	mock.ExpectCommit()

	result, err := repo.CreateTransactionWithOptions(context.Background(), items, opts)
	assert.NoError(t, err)
	assert.Equal(t, 10000.0, result.TotalAmount)
	assert.NoError(t, mock.ExpectationsWereMet())
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./internal/repositories -v -run "TestTransactionRepository_CreateTransactionWithOptions_(Default|Batch|WithRow)"
```

Expected: FAIL - tests will fail because stub implementation doesn't handle options correctly

**Step 3: Implement unified CreateTransactionWithOptions**

Replace the stub implementation with the full unified logic:

```go
// Replace stub in internal/repositories/transaction_repository.go
func (repo *transactionRepository) CreateTransactionWithOptions(ctx context.Context, items []models.CheckoutItem, opts TransactionOptions) (*models.Transaction, error) {
	// Begin transaction with appropriate isolation level
	txOptions := &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	tx, err := repo.db.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var totalAmount float64
	var details []models.TransactionDetail

	// Branch based on batch operations flag
	if opts.UseBatchOperations {
		totalAmount, details, err = repo.processBatchOperations(ctx, tx, items, opts.UseRowLocking)
	} else {
		totalAmount, details, err = repo.processLoopOperations(ctx, tx, items)
	}
	if err != nil {
		return nil, err
	}

	// Insert transaction
	var transactionID int
	var createdAt time.Time
	var updatedAt *time.Time
	err = tx.QueryRowContext(ctx,
		`INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at, updated_at`,
		totalAmount,
	).Scan(&transactionID, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert transaction: %w", err)
	}

	// Insert transaction details
	if opts.UseBatchOperations {
		details, err = repo.insertDetailsBatch(ctx, tx, transactionID, details)
	} else {
		details, err = repo.insertDetailsLoop(ctx, tx, transactionID, details)
	}
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

// Helper: Process products using loop-based queries
func (repo *transactionRepository) processLoopOperations(ctx context.Context, tx *sql.Tx, items []models.CheckoutItem) (float64, []models.TransactionDetail, error) {
	var totalAmount float64
	details := make([]models.TransactionDetail, 0, len(items))

	for _, item := range items {
		var id int
		var name string
		var price float64
		var stock int

		err := tx.QueryRowContext(ctx,
			`SELECT id, name, price, stock FROM products WHERE id=$1`,
			item.ProductID,
		).Scan(&id, &name, &price, &stock)
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, nil, errors.New("product not found")
			}
			return 0, nil, fmt.Errorf("get product: %w", err)
		}

		if stock < item.Quantity {
			return 0, nil, fmt.Errorf("insufficient stock for product %d", item.ProductID)
		}

		subtotal := float64(item.Quantity) * price
		totalAmount += subtotal

		_, err = tx.ExecContext(ctx,
			`UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2`,
			item.Quantity, id,
		)
		if err != nil {
			return 0, nil, fmt.Errorf("update product stock: %w", err)
		}

		details = append(details, models.TransactionDetail{
			ProductID:   id,
			ProductName: name,
			Price:       price,
			Quantity:    item.Quantity,
			SubTotal:    subtotal,
		})
	}

	return totalAmount, details, nil
}

// Helper: Process products using batch queries
func (repo *transactionRepository) processBatchOperations(ctx context.Context, tx *sql.Tx, items []models.CheckoutItem, useRowLocking bool) (float64, []models.TransactionDetail, error) {
	productIDs := make([]int, 0, len(items))
	itemQtyMap := make(map[int]int, len(items))

	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
		itemQtyMap[item.ProductID] = item.Quantity
	}

	// Build SELECT query
	placeholders := make([]string, len(productIDs))
	args := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(
		`SELECT id, name, price, stock FROM products WHERE id IN (%s)`,
		strings.Join(placeholders, ", "),
	)
	if useRowLocking {
		query += " FOR UPDATE"
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, nil, fmt.Errorf("select products: %w", err)
	}
	defer rows.Close()

	productMap := make(map[int]struct {
		Name  string
		Price float64
		Stock int
	}, len(productIDs))

	for rows.Next() {
		var id int
		var name string
		var price float64
		var stock int
		if err := rows.Scan(&id, &name, &price, &stock); err != nil {
			return 0, nil, fmt.Errorf("scan product: %w", err)
		}
		productMap[id] = struct {
			Name  string
			Price float64
			Stock int
		}{name, price, stock}
	}
	if err := rows.Err(); err != nil {
		return 0, nil, fmt.Errorf("iterate products: %w", err)
	}

	if len(productMap) != len(productIDs) {
		return 0, nil, errors.New("one or more products not found")
	}

	// Validate stock and calculate totals
	var totalAmount float64
	details := make([]models.TransactionDetail, 0, len(items))

	for _, item := range items {
		product := productMap[item.ProductID]
		if product.Stock < item.Quantity {
			return 0, nil, fmt.Errorf("insufficient stock for product %d", item.ProductID)
		}

		subtotal := float64(item.Quantity) * product.Price
		totalAmount += subtotal

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
			SubTotal:    subtotal,
		})
	}

	// Batch update stocks
	updateIDs := make([]int, 0, len(items))
	updateQtys := make([]int, 0, len(items))
	for _, item := range items {
		updateIDs = append(updateIDs, item.ProductID)
		updateQtys = append(updateQtys, item.Quantity)
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE products AS p SET stock = p.stock - u.qty, updated_at = NOW()
		 FROM (SELECT UNNEST($1::int[]) AS id, UNNEST($2::int[]) AS qty) AS u
		 WHERE p.id = u.id`,
		pq.Array(updateIDs), pq.Array(updateQtys),
	)
	if err != nil {
		return 0, nil, fmt.Errorf("batch update stocks: %w", err)
	}

	return totalAmount, details, nil
}

// Helper: Insert details using loop
func (repo *transactionRepository) insertDetailsLoop(ctx context.Context, tx *sql.Tx, transactionID int, details []models.TransactionDetail) ([]models.TransactionDetail, error) {
	for i := range details {
		details[i].TransactionID = transactionID

		var id int
		var createdAt time.Time
		var updatedAt *time.Time

		err := tx.QueryRowContext(ctx,
			`INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id, created_at, updated_at`,
			transactionID, details[i].ProductID, details[i].Price, details[i].Quantity, details[i].SubTotal,
		).Scan(&id, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("insert transaction detail: %w", err)
		}

		details[i].ID = id
		details[i].CreatedAt = createdAt
		details[i].UpdatedAt = updatedAt
	}
	return details, nil
}

// Helper: Insert details using batch
func (repo *transactionRepository) insertDetailsBatch(ctx context.Context, tx *sql.Tx, transactionID int, details []models.TransactionDetail) ([]models.TransactionDetail, error) {
	if len(details) == 0 {
		return details, nil
	}

	txIDs := make([]int, len(details))
	productIDs := make([]int, len(details))
	prices := make([]float64, len(details))
	qtys := make([]int, len(details))
	subtotals := make([]float64, len(details))

	for i, d := range details {
		txIDs[i] = transactionID
		productIDs[i] = d.ProductID
		prices[i] = d.Price
		qtys[i] = d.Quantity
		subtotals[i] = d.SubTotal
	}

	rows, err := tx.QueryContext(ctx,
		`INSERT INTO transaction_details (transaction_id, product_id, price, quantity, subtotal)
		 SELECT * FROM UNNEST($1::int[], $2::int[], $3::float8[], $4::int[], $5::float8[])
		 RETURNING id, transaction_id, product_id, price, quantity, subtotal, created_at, updated_at`,
		pq.Array(txIDs), pq.Array(productIDs), pq.Array(prices), pq.Array(qtys), pq.Array(subtotals),
	)
	if err != nil {
		return nil, fmt.Errorf("batch insert details: %w", err)
	}
	defer rows.Close()

	// Reconstruct details with DB values
	productNameMap := make(map[int]string, len(details))
	for _, d := range details {
		productNameMap[d.ProductID] = d.ProductName
	}

	result := make([]models.TransactionDetail, 0, len(details))
	for rows.Next() {
		var d models.TransactionDetail
		err := rows.Scan(&d.ID, &d.TransactionID, &d.ProductID, &d.Price, &d.Quantity, &d.SubTotal, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan detail: %w", err)
		}
		d.ProductName = productNameMap[d.ProductID]
		result = append(result, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate details: %w", err)
	}

	return result, nil
}
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/repositories -v -run "TestTransactionRepository_CreateTransactionWithOptions"
```

Expected: PASS (all new tests + interface test)

**Step 5: Commit**

```bash
git add internal/repositories/transaction_repository.go internal/repositories/transaction_repository_test.go
git commit -m "feat: implement unified CreateTransactionWithOptions with configurable batch/lock options"
```

---

## Task 4: Update Service Layer

**Files:**
- Modify: `internal/services/transaction_service.go:22-25`
- Test: `internal/services/transaction_service_test.go`

**Step 1: Write failing test for service layer**

```go
// Add to internal/services/transaction_service_test.go
func TestTransactionService_Checkout_UseLockTrue(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)
	service := NewTrasactionService(mockRepo)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
	}

	expectedTx := &models.Transaction{
		ID:          1,
		TotalAmount: 20000.0,
	}

	// Expect CreateTransactionWithOptions to be called with optimal settings
	mockRepo.On("CreateTransactionWithOptions", 
		mock.Anything, 
		items, 
		repositories.OptimalTransactionOptions(),
	).Return(expectedTx, nil)

	result, err := service.Checkout(context.Background(), items, true)

	assert.NoError(t, err)
	assert.Equal(t, expectedTx, result)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_Checkout_UseLockFalse(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)
	service := NewTrasactionService(mockRepo)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
	}

	expectedTx := &models.Transaction{
		ID:          1,
		TotalAmount: 20000.0,
	}

	// Expect CreateTransactionWithOptions with default settings
	mockRepo.On("CreateTransactionWithOptions",
		mock.Anything,
		items,
		repositories.DefaultTransactionOptions(),
	).Return(expectedTx, nil)

	result, err := service.Checkout(context.Background(), items, false)

	assert.NoError(t, err)
	assert.Equal(t, expectedTx, result)
	mockRepo.AssertExpectations(t)
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./internal/services -v -run TestTransactionService_Checkout_UseLock
```

Expected: FAIL - mock expectations not met (still calling CreateTransaction)

**Step 3: Update service implementation**

```go
// Modify internal/services/transaction_service.go:22-25
func (serv *transactionService) Checkout(ctx context.Context, items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	var opts repositories.TransactionOptions
	
	if useLock {
		opts = repositories.OptimalTransactionOptions()
	} else {
		opts = repositories.DefaultTransactionOptions()
	}
	
	return serv.transactionRepo.CreateTransactionWithOptions(ctx, items, opts)
}
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/services -v -run TestTransactionService_Checkout_UseLock
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/services/transaction_service.go internal/services/transaction_service_test.go
git commit -m "feat: update service to use CreateTransactionWithOptions based on useLock flag"
```

---

## Task 5: Update Mock Repository

**Files:**
- Modify: `internal/mocks/transaction_repository_mock.go`

**Step 1: Add new method to mock**

```go
// Add to internal/mocks/transaction_repository_mock.go
func (m *MockTransactionRepository) CreateTransactionWithOptions(ctx context.Context, items []models.CheckoutItem, opts repositories.TransactionOptions) (*models.Transaction, error) {
	args := m.Called(ctx, items, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}
```

**Step 2: Run all service tests to verify mock works**

```bash
go test ./internal/services -v
```

Expected: PASS (all tests)

**Step 3: Commit**

```bash
git add internal/mocks/transaction_repository_mock.go
git commit -m "feat: add CreateTransactionWithOptions to mock repository"
```

---

## Task 6: Deprecate Old Methods (Optional - Keep for Backward Compatibility)

**Files:**
- Modify: `internal/repositories/transaction_repository.go:20-30,35-270`

**Step 1: Add deprecation comments**

```go
// Modify internal/repositories/transaction_repository.go

// CreateTransaction creates a transaction using loop-based operations.
// Deprecated: Use CreateTransactionWithOptions with DefaultTransactionOptions() instead.
func (repo *transactionRepository) CreateTransaction(ctx context.Context, items []models.CheckoutItem) (*models.Transaction, error) {
	return repo.CreateTransactionWithOptions(ctx, items, DefaultTransactionOptions())
}

// CreateTransactionOptimal creates a transaction using batch operations with row locking.
// Deprecated: Use CreateTransactionWithOptions with OptimalTransactionOptions() instead.
func (repo *transactionRepository) CreateTransactionOptimal(ctx context.Context, items []models.CheckoutItem) (*models.Transaction, error) {
	return repo.CreateTransactionWithOptions(ctx, items, OptimalTransactionOptions())
}
```

**Step 2: Update existing implementations to delegate**

Replace the full implementations of `CreateTransaction` and `CreateTransactionOptimal` with simple delegations (as shown above).

**Step 3: Run all repository tests**

```bash
go test ./internal/repositories -v
```

Expected: PASS (all tests, including old ones)

**Step 4: Commit**

```bash
git add internal/repositories/transaction_repository.go
git commit -m "refactor: deprecate old methods in favor of CreateTransactionWithOptions"
```

---

## Task 7: Run Full Test Suite

**Files:**
- N/A (verification only)

**Step 1: Run all tests with coverage**

```bash
cd /home/fajarsiagian/Learning/golang/go-kasir-umam-ds
go test ./... -coverprofile=coverage.out
```

Expected: PASS with coverage ≥ 87%

**Step 2: Check coverage report**

```bash
go tool cover -func=coverage.out | grep total
```

Expected: Total coverage should be at least 87.3%

**Step 3: Run specific package tests**

```bash
go test ./internal/repositories -v
go test ./internal/services -v
go test ./handlers -v
```

Expected: All PASS

**Step 4: Verify no regressions**

```bash
make test
make check-coverage
```

Expected: All commands succeed

---

## Task 8: Update Documentation

**Files:**
- Create: `docs/architecture/transaction-repository-design.md`
- Modify: `docs/README.md`

**Step 1: Create architecture documentation**

```markdown
<!-- docs/architecture/transaction-repository-design.md -->
# Transaction Repository Architecture

## Overview

The transaction repository provides a unified interface for creating transactions with configurable performance and concurrency characteristics.

## Design Philosophy

Instead of having separate methods for different operation modes, we use a single configurable method that adapts based on runtime options.

## TransactionOptions

```go
type TransactionOptions struct {
    UseBatchOperations bool  // Batch INSERT/UPDATE with UNNEST
    UseRowLocking      bool  // SELECT ... FOR UPDATE
}
```

**Presets:**
- `DefaultTransactionOptions()`: Conservative, loop-based (no batch, no locking)
- `OptimalTransactionOptions()`: High-performance (batch + locking)

## Architecture

```
Service Layer
    ↓
transactionService.Checkout(useLock) 
    ↓ (maps useLock → TransactionOptions)
    ↓
Repository Layer
    ↓
CreateTransactionWithOptions(opts)
    ↓
    ├─ opts.UseBatchOperations?
    │   ├─ YES → processBatchOperations()
    │   └─ NO  → processLoopOperations()
    ↓
    ├─ Insert Transaction
    ↓
    ├─ opts.UseBatchOperations?
    │   ├─ YES → insertDetailsBatch()
    │   └─ NO  → insertDetailsLoop()
    ↓
    └─ Commit & Return
```

## Implementation Strategies

### Loop-Based (Default)
- Individual SELECT for each product
- Individual UPDATE for each stock change
- Individual INSERT for each transaction detail
- **Best for:** Low-volume, simple transactions
- **Tradeoff:** More DB round-trips, simpler code

### Batch-Based (Optimal)
- Single SELECT IN with all product IDs
- Single UPDATE with UNNEST for all stock changes
- Single INSERT with UNNEST for all details
- **Best for:** High-volume, multi-item transactions
- **Tradeoff:** Fewer DB round-trips, more complex queries

### Row Locking
- Adds `FOR UPDATE` to SELECT queries
- Prevents concurrent stock modifications
- **Best for:** High-concurrency scenarios
- **Tradeoff:** Potential lock contention

## Usage Examples

```go
// Conservative (default)
opts := repositories.DefaultTransactionOptions()
tx, err := repo.CreateTransactionWithOptions(ctx, items, opts)

// High-performance
opts := repositories.OptimalTransactionOptions()
tx, err := repo.CreateTransactionWithOptions(ctx, items, opts)

// Custom configuration
opts := repositories.TransactionOptions{
    UseBatchOperations: true,
    UseRowLocking:      false,
}
tx, err := repo.CreateTransactionWithOptions(ctx, items, opts)
```

## Service Layer Integration

The service layer maps the boolean `useLock` parameter to appropriate options:

```go
func (s *transactionService) Checkout(ctx context.Context, items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
    var opts repositories.TransactionOptions
    if useLock {
        opts = repositories.OptimalTransactionOptions()
    } else {
        opts = repositories.DefaultTransactionOptions()
    }
    return s.transactionRepo.CreateTransactionWithOptions(ctx, items, opts)
}
```

## Backward Compatibility

The original methods are preserved but deprecated:

```go
CreateTransaction()        → CreateTransactionWithOptions(DefaultTransactionOptions())
CreateTransactionOptimal() → CreateTransactionWithOptions(OptimalTransactionOptions())
```

## Performance Characteristics

| Configuration | DB Queries (10 items) | Concurrency Safety | Use Case |
|---------------|----------------------|-------------------|----------|
| Default | ~31 queries | Moderate | Development, low-traffic |
| Batch only | ~4 queries | Moderate | High-traffic, low contention |
| Batch + Lock | ~4 queries | High | High-traffic, high contention |

## Testing Strategy

- Unit tests cover all option combinations
- Mocks verify correct option propagation
- Integration tests validate SQL generation
- Coverage maintained at ≥87%

## Future Enhancements

- Add `context.Context` deadline handling
- Implement retry logic for deadlocks
- Add metrics/tracing integration
- Support for distributed transactions
```

**Step 2: Update docs/README.md**

```markdown
<!-- Add to docs/README.md under Architecture section -->

## Architecture Documentation

- [Overview](architecture/overview.md) - High-level system architecture
- [Clean Architecture](architecture/clean-architecture.md) - Clean architecture patterns
- [Database Schema](architecture/database-schema.md) - Database design
- **[Transaction Repository Design](architecture/transaction-repository-design.md)** - Configurable transaction creation ✨ NEW
```

**Step 3: Commit documentation**

```bash
git add docs/architecture/transaction-repository-design.md docs/README.md
git commit -m "docs: add transaction repository architecture documentation"
```

---

## Task 9: Integration Testing (Optional)

**Files:**
- Create: `internal/repositories/transaction_repository_integration_test.go`

**Step 1: Create integration test**

```go
// internal/repositories/transaction_repository_integration_test.go
//go:build integration

package repositories

import (
	"context"
	"database/sql"
	"fajar7xx/go-kasir-umam-ds/models"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: This test requires a real PostgreSQL database
// Set DB_CONN environment variable: export DB_CONN="postgres://user:pass@localhost/testdb?sslmode=disable"

func setupTestDB(t *testing.T) *sql.DB {
	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		t.Skip("DB_CONN not set, skipping integration test")
	}

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	// Create test tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			price FLOAT8,
			stock INT,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			total_amount FLOAT8,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS transaction_details (
			id SERIAL PRIMARY KEY,
			transaction_id INT REFERENCES transactions(id),
			product_id INT REFERENCES products(id),
			price FLOAT8,
			quantity INT,
			subtotal FLOAT8,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP
		);
	`)
	require.NoError(t, err)

	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`
		TRUNCATE transaction_details, transactions, products RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
	db.Close()
}

func TestIntegration_CreateTransactionWithOptions_AllModes(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewTransactionRepository(db)

	// Insert test products
	_, err := db.Exec(`
		INSERT INTO products (name, price, stock) VALUES
		('Product A', 10000, 100),
		('Product B', 20000, 50)
	`)
	require.NoError(t, err)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
		{ProductID: 2, Quantity: 1},
	}

	testCases := []struct {
		name string
		opts TransactionOptions
	}{
		{"Default", DefaultTransactionOptions()},
		{"Batch", TransactionOptions{UseBatchOperations: true}},
		{"Lock", TransactionOptions{UseRowLocking: true}},
		{"Optimal", OptimalTransactionOptions()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean between runs
			_, _ = db.Exec(`TRUNCATE transaction_details, transactions RESTART IDENTITY CASCADE`)

			result, err := repo.CreateTransactionWithOptions(context.Background(), items, tc.opts)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, 40000.0, result.TotalAmount) // 2*10000 + 1*20000
			assert.Len(t, result.Details, 2)
			assert.NotZero(t, result.ID)
			assert.NotZero(t, result.CreatedAt)

			// Verify stock was updated
			var stockA, stockB int
			err = db.QueryRow(`SELECT stock FROM products WHERE id = 1`).Scan(&stockA)
			require.NoError(t, err)
			assert.Equal(t, 98, stockA) // 100 - 2

			err = db.QueryRow(`SELECT stock FROM products WHERE id = 2`).Scan(&stockB)
			require.NoError(t, err)
			assert.Equal(t, 49, stockB) // 50 - 1
		})
	}
}
```

**Step 2: Run integration test (if DB available)**

```bash
export DB_CONN="postgres://user:pass@localhost/testdb?sslmode=disable"
go test ./internal/repositories -tags=integration -v
```

Expected: PASS (if DB configured), SKIP (if not)

**Step 3: Commit integration test**

```bash
git add internal/repositories/transaction_repository_integration_test.go
git commit -m "test: add integration tests for transaction repository"
```

---

## Completion Checklist

- [ ] Task 1: TransactionOptions type created and tested
- [ ] Task 2: Repository interface updated
- [ ] Task 3: Unified CreateTransactionWithOptions implemented
- [ ] Task 4: Service layer updated to use new method
- [ ] Task 5: Mock repository updated
- [ ] Task 6: Old methods deprecated
- [ ] Task 7: Full test suite passes
- [ ] Task 8: Documentation updated
- [ ] Task 9: Integration tests added (optional)

## Success Criteria

✅ All tests passing with ≥87% coverage  
✅ Service layer correctly uses `useLock` parameter  
✅ Single source of truth for transaction creation  
✅ Backward compatibility maintained  
✅ Documentation updated

## Rollback Plan

If issues arise during implementation:

1. Revert last commit: `git reset --hard HEAD~1`
2. Check test status: `go test ./...`
3. Return to previous stable state

---

**End of Implementation Plan**

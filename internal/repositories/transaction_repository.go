package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fajar7xx/go-kasir-umam-ds/models"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, items []models.CheckoutItem) (*models.Transaction, error)
	CreateTransactionOptimal(ctx context.Context, items []models.CheckoutItem) (*models.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (repo *transactionRepository) CreateTransaction(ctx context.Context, items []models.CheckoutItem) (*models.Transaction, error) {
	// database transaction
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error begin transaction %w:", err)
	}
	defer tx.Rollback()

	// init subtotal
	totalAmount := 0.0
	// init transaction_details
	// details := make([]models.TransactionDetail, 0)
	// FIX 1: Pre-allocate capacity, tapi length = 0
	details := make([]models.TransactionDetail, 0, len(items))

	// loop setiap item
	for _, item := range items {
		var productName string
		var productPrice float64
		var productID, productStock int

		// get pricing each product
		getProductQuery := `SELECT
							id,
							name,
							price,
							stock
							FROM
							products
							WHERE id=$1`
		err := tx.QueryRowContext(ctx, getProductQuery, item.ProductID).Scan(
			&productID,
			&productName,
			&productPrice,
			&productStock,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("product not found")
			}
			return nil, fmt.Errorf("error to get product: %w", err)
		}

		// multiple each item with qty
		// add to subtotal
		subtotal := float64(item.Quantity) * productPrice
		// calculate totalAMount
		totalAmount += subtotal

		// kurangin jumlah stock
		updateProductQuery := `UPDATE products
								SET
								stock = stock - $1,
								updated_at = NOW()
								WHERE id = $2`
		_, err = tx.ExecContext(ctx, updateProductQuery,
			item.Quantity,
			productID,
		)
		if err != nil {
			return nil, fmt.Errorf("error to update product stock: %w", err)
		}

		// item dimasukkan ke transactionDetails
		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Price:       productPrice,
			Quantity:    item.Quantity,
			SubTotal:    subtotal,
		})
	}

	// insert transactions
	var transactionID int
	var transactionCreatedAt time.Time
	var transactionUpdatedAt *time.Time
	insertTransactionQuery := `INSERT INTO transactions
								(total_amount)
								VALUES
								($1)
								RETURNING id, created_at, updated_at`
	err = tx.QueryRowContext(ctx, insertTransactionQuery, totalAmount).Scan(
		&transactionID,
		&transactionCreatedAt,
		&transactionUpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error to insert transaction: %w", err)
	}

	// insert transaction_details
	insertTransactionDetailQuery := `INSERT INTO transaction_details
									(transaction_id, product_id, price, quantity, subtotal)
									VALUES
									($1, $2, $3, $4, $5)
									RETURNING id, created_at, updated_at`
	for i := range details {
		details[i].TransactionID = transactionID

		var detailID int
		var detailCreatedAt time.Time
		var detailUpdatedAt *time.Time

		err = tx.QueryRowContext(ctx, insertTransactionDetailQuery,
			details[i].TransactionID,
			details[i].ProductID,
			details[i].Price,
			details[i].Quantity,
			details[i].SubTotal,
		).Scan(&detailID, &detailCreatedAt, &detailUpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error to insert transaction detail: %w", err)
		}

		// Update detail with database values
		details[i].ID = detailID
		details[i].CreatedAt = detailCreatedAt
		details[i].UpdatedAt = detailUpdatedAt
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error to commit transaction: %w", err)
	}

	result := &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
		CreatedAt:   transactionCreatedAt,
		UpdatedAt:   transactionUpdatedAt,
	}

	return result, nil
}

func (repo *transactionRepository) CreateTransactionOptimal(ctx context.Context, items []models.CheckoutItem) (*models.Transaction, error) {
	// database transaction
	//
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, fmt.Errorf("error begin transaction %w:", err)
	}
	defer tx.Rollback()

	// ==========================================
	// FIX #1: BATCH SELECT dengan IN clause + FOR UPDATE (row locking)
	// ==========================================
	productIDs := make([]int, 0, len(items))
	itemQtyMap := make(map[int]int, len(items)) //productID => quantity

	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
		itemQtyMap[item.ProductID] = item.Quantity
	}

	// build placeholders
	placeholders := make([]string, len(productIDs))
	args := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	// Single SELECT dengan FOR UPDATE untuk lock rows
	query := fmt.Sprintf(`SELECT id, name, price, stock
							FROM products
							WHERE id IN (%s)
							FOR UPDATE`,
		strings.Join(placeholders, ", "))
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error selecting products %w:", err)
	}
	defer rows.Close()

	// Map products untuk O(1) lookup
	productMap := make(map[int]struct {
		Name  string
		Price float64
		Stock int
	}, len(productIDs))

	for rows.Next() {
		var ID int
		var name string
		var price float64
		var stock int
		if err := rows.Scan(&ID, &name, &price, &stock); err != nil {
			return nil, fmt.Errorf("error scanning product %w:", err)
		}
		productMap[ID] = struct {
			Name  string
			Price float64
			Stock int
		}{name, price, stock}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}

	// Validasi: semua product harus ada
	if len(productMap) != len(productIDs) {
		return nil, errors.New("one or more products not found")
	}

	// ==========================================
	// FIX #2: VALIDASI STOCK sebelum processing
	// ==========================================
	var totalAmount float64
	details := make([]models.TransactionDetail, 0, len(items))

	for _, item := range items {
		product, exists := productMap[item.ProductID]
		if !exists {
			return nil, fmt.Errorf("product %d not found", item.ProductID)
		}

		// Check stock availability
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %d: requested %d, available %d",
				item.ProductID, item.Quantity, product.Stock)
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

	// ==========================================
	// FIX #3: BATCH UPDATE menggunakan UNNEST
	// ==========================================
	// Construct arrays untuk batch update
	updateIDs := make([]int, 0, len(items))
	updateQtys := make([]int, 0, len(items))

	for _, item := range items {
		updateIDs = append(updateIDs, item.ProductID)
		updateQtys = append(updateQtys, item.Quantity)
	}

	updateQuery := `
			UPDATE products AS p
			SET
				stock = p.stock - u.qty,
				updated_at = NOW()
			FROM (
				SELECT
					UNNEST($1::int[]) AS id,
					UNNEST($2::int[]) AS qty
			) AS u
			WHERE p.id = u.id
		`

	_, err = tx.ExecContext(ctx, updateQuery, pq.Array(updateIDs), pq.Array(updateQtys))
	if err != nil {
		return nil, fmt.Errorf("update product stock: %w", err)
	}

	// ==========================================
	// FIX #4: INSERT transaction
	// ==========================================
	var transactionID int
	var transactionCreatedAt time.Time
	var transactionUpdatedAt *time.Time
	insertTxQuery := `
			INSERT INTO transactions (total_amount)
			VALUES ($1)
			RETURNING id, created_at, updated_at
		`
	err = tx.QueryRowContext(ctx, insertTxQuery, totalAmount).Scan(
		&transactionID,
		&transactionCreatedAt,
		&transactionUpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert transaction: %w", err)
	}

	// ==========================================
	// FIX #5: BATCH INSERT transaction details menggunakan UNNEST
	// ==========================================
	if len(details) > 0 {
		detailTxIDs := make([]int, len(details))
		detailProductIDs := make([]int, len(details))
		detailPrices := make([]float64, len(details))
		detailQtys := make([]int, len(details))
		detailSubtotals := make([]float64, len(details))

		for i, detail := range details {
			detailTxIDs[i] = transactionID
			detailProductIDs[i] = detail.ProductID
			detailPrices[i] = detail.Price
			detailQtys[i] = detail.Quantity
			detailSubtotals[i] = detail.SubTotal
		}

		insertDetailsQuery := `
				INSERT INTO transaction_details
					(transaction_id, product_id, price, quantity, subtotal)
				SELECT * FROM UNNEST(
					$1::int[],
					$2::int[],
					$3::float8[],
					$4::int[],
					$5::float8[]
				)
				RETURNING id, transaction_id, product_id, price, quantity, subtotal, created_at, updated_at
			`

		rows, err := tx.QueryContext(ctx, insertDetailsQuery,
			pq.Array(detailTxIDs),
			pq.Array(detailProductIDs),
			pq.Array(detailPrices),
			pq.Array(detailQtys),
			pq.Array(detailSubtotals),
		)
		if err != nil {
			return nil, fmt.Errorf("insert transaction details: %w", err)
		}
		defer rows.Close()

		// Scan results and reconstruct with all database values
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
				return nil, fmt.Errorf("scan transaction detail: %w", err)
			}

			// Add product name from existing productMap (O(1) lookup)
			if product, exists := productMap[d.ProductID]; exists {
				d.ProductName = product.Name
			}

			finalDetails = append(finalDetails, d)
		}

		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("iterate transaction details: %w", err)
		}

		// Replace details with finalDetails containing database values
		details = finalDetails
	}

	// ==========================================
	// COMMIT
	// ==========================================
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
		CreatedAt:   transactionCreatedAt,
		UpdatedAt:   transactionUpdatedAt,
	}, nil
}

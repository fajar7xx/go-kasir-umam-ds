package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fajar7xx/go-kasir-umam-ds/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestProductRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	now := time.Now()
	desc := "Delicious Food"
	categoryDesc := "Food Category"

	// Updated query to match actual implementation with JOIN
	query := regexp.QuoteMeta(`SELECT
				  p.id,
				  p.name,
				  p.description,
				  p.price,
				  p.stock,
				  p.category_id,
				  p.created_at,
				  p.updated_at,
				  c.id,
				  c.name,
				  c.description
				FROM
				  products p
				  JOIN categories c ON p.category_id = c.id
				ORDER BY p.created_at DESC`)

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock", "category_id", "created_at", "updated_at",
		"category_id", "category_name", "category_description",
	}).
		AddRow(1, "Nasi Goreng", &desc, 15000.0, 10, 1, now, now, 1, "Food", &categoryDesc).
		AddRow(2, "Es Teh", nil, 3000.0, 20, 2, now, now, 2, "Beverage", nil)

	mock.ExpectQuery(query).WillReturnRows(rows)

	products, err := repo.GetAll(context.Background(), "")

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Nasi Goreng", products[0].Name)
	assert.NotNil(t, products[0].Description)
	assert.Equal(t, "Delicious Food", *products[0].Description)
	assert.Equal(t, 1, products[0].Category.ID)
	assert.Equal(t, "Food", products[0].Category.Name)
	assert.Equal(t, "Es Teh", products[1].Name)
	assert.Nil(t, products[1].Description)
	assert.Equal(t, 2, products[1].Category.ID)
}

func TestProductRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	now := time.Now()
	desc := "Delicious Food"
	categoryDesc := "Food Category"

	query := regexp.QuoteMeta(`select
				  p.id,
				  p.name,
				  p.description,
				  p.price,
				  p.stock,
				  p.category_id,
				  p.created_at,
				  p.updated_at,
				  c.id as category_id,
				  c.name as category_name,
				  c.description as category_description
				from
				  products p
				  join categories c on p.category_id = c.id
				where p.id = $1`)

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock", "category_id", "created_at", "updated_at",
		"category_id", "category_name", "category_description",
	}).AddRow(1, "Nasi Goreng", &desc, 15000.0, 10, 1, now, now, 1, "Food", &categoryDesc)

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	product, err := repo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, 1, product.ID)
	assert.Equal(t, "Nasi Goreng", product.Name)
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	query := regexp.QuoteMeta(`select
				  p.id,
				  p.name,
				  p.description,
				  p.price,
				  p.stock,
				  p.category_id,
				  p.created_at,
				  p.updated_at,
				  c.id as category_id,
				  c.name as category_name,
				  c.description as category_description
				from
				  products p
				  join categories c on p.category_id = c.id
				where p.id = $1`)
	mock.ExpectQuery(query).WithArgs(1).WillReturnError(sql.ErrNoRows)

	product, err := repo.GetByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
	assert.Nil(t, product)
}

func TestProductRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	now := time.Now()
	desc := "Delicious Food"
	categoryDesc := "Food Category"
	product := &models.Product{
		Name:        "Nasi Goreng",
		Description: &desc,
		Price:       15000.0,
		Stock:       10,
		CategoryID:  1,
	}

	query := regexp.QuoteMeta(`INSERT INTO products (name, price, stock, description, category_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`)

	mock.ExpectQuery(query).
		WithArgs(product.Name, product.Price, product.Stock, product.Description, product.CategoryID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, now, now))

	// Mock GetByID call after Create
	getQuery := regexp.QuoteMeta(`select
				  p.id,
				  p.name,
				  p.description,
				  p.price,
				  p.stock,
				  p.category_id,
				  p.created_at,
				  p.updated_at,
				  c.id as category_id,
				  c.name as category_name,
				  c.description as category_description
				from
				  products p
				  join categories c on p.category_id = c.id
				where p.id = $1`)
	mock.ExpectQuery(getQuery).WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{
			"id", "name", "description", "price", "stock", "category_id", "created_at", "updated_at",
			"category_id", "category_name", "category_description",
		}).AddRow(1, "Nasi Goreng", &desc, 15000.0, 10, 1, now, now, 1, "Food", &categoryDesc),
	)

	result, err := repo.Create(context.Background(), product)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Nasi Goreng", result.Name)
	assert.Equal(t, "Food", result.Category.Name)
}

func TestProductRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	desc := "Updated Desc"
	product := &models.Product{
		Name:        "Nasi Goreng Updated",
		Description: &desc,
		Price:       16000.0,
		Stock:       15,
		CategoryID:  1,
	}

	query := regexp.QuoteMeta(`UPDATE products SET name = $1, price=$2, stock=$3, description=$4, category_id=$5, updated_at = NOW() WHERE id = $6`)
	mock.ExpectExec(query).
		WithArgs(product.Name, product.Price, product.Stock, product.Description, product.CategoryID, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), 1, product)

	assert.NoError(t, err)
}

func TestProductRepository_Update_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	desc := "Updated Desc"
	product := &models.Product{
		Name:        "Nasi Goreng Updated",
		Description: &desc,
		Price:       16000.0,
		Stock:       15,
		CategoryID:  1,
	}

	query := regexp.QuoteMeta(`UPDATE products SET name = $1, price=$2, stock=$3, description=$4, category_id=$5, updated_at = NOW() WHERE id = $6`)
	mock.ExpectExec(query).
		WithArgs(product.Name, product.Price, product.Stock, product.Description, product.CategoryID, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Update(context.Background(), 1, product)

	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
}

func TestProductRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	query := regexp.QuoteMeta(`DELETE from products where id = $1`)
	mock.ExpectExec(query).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), 1)

	assert.NoError(t, err)
}

func TestProductRepository_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	query := regexp.QuoteMeta(`DELETE from products where id = $1`)
	mock.ExpectExec(query).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(context.Background(), 1)

	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
}

func TestProductRepository_GetAll_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	now := time.Now()
	desc := "Delicious product"
	categoryDesc := "Food Category"

	// Return rows but inject error during iteration - must match all columns from the actual query
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "category_id", "created_at", "updated_at", "category_id", "category_name", "category_description"}).
		AddRow(1, "Burger", &desc, 25000.0, 10, 1, now, now, 1, "Food", &categoryDesc).
		RowError(0, errors.New("row iteration error"))

	// Match the actual query from product_repository.go (lines 36-47)
	query := regexp.QuoteMeta(`SELECT
				  p.id,
				  p.name,
				  p.description,
				  p.price,
				  p.stock,
				  p.category_id,
				  p.created_at,
				  p.updated_at,
				  c.id,
				  c.name,
				  c.description
				FROM
				  products p
				  JOIN categories c ON p.category_id = c.id
				ORDER BY p.created_at DESC`)
	mock.ExpectQuery(query).WillReturnRows(rows)

	products, err := repo.GetAll(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, products)
	assert.Contains(t, err.Error(), "row iteration error")
}

func TestProductRepository_Update_RowsAffectedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	desc := "Product Updated"
	product := &models.Product{
		Name:        "Product Updated",
		Description: &desc,
		Price:       20000,
		Stock:       15,
		CategoryID:  1,
	}

	query := regexp.QuoteMeta(`UPDATE products SET name = $1, price=$2, stock=$3, description=$4, category_id=$5, updated_at = NOW() WHERE id = $6`)

	// Return error on RowsAffected()
	mock.ExpectExec(query).
		WithArgs(product.Name, product.Price, product.Stock, product.Description, product.CategoryID, 1).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

	err = repo.Update(context.Background(), 1, product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rows affected error")
}

func TestProductRepository_Delete_RowsAffectedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	query := regexp.QuoteMeta(`DELETE from products where id = $1`)

	mock.ExpectExec(query).
		WithArgs(1).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

	err = repo.Delete(context.Background(), 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rows affected error")
}

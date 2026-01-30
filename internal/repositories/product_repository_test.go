package repositories

import (
	"database/sql"
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
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "category_id", "created_at", "updated_at"}).
		AddRow(1, "Nasi Goreng", &desc, 15000.0, 10, 1, now, now).
		AddRow(2, "Es Teh", nil, 3000.0, 20, 2, now, now)

	query := regexp.QuoteMeta(`SELECT id, name, description, price, stock, category_id, created_at, updated_at FROM products`)
	mock.ExpectQuery(query).WillReturnRows(rows)

	products, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Nasi Goreng", products[0].Name)
	assert.NotNil(t, products[0].Description)
	assert.Equal(t, "Delicious Food", *products[0].Description)
	assert.Equal(t, "Es Teh", products[1].Name)
	assert.Nil(t, products[1].Description)
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
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "category_id", "created_at", "updated_at"}).
		AddRow(1, "Nasi Goreng", &desc, 15000.0, 10, 1, now, now)

	query := regexp.QuoteMeta(`SELECT id, name, description, price, stock, category_id, created_at, updated_at FROM products WHERE id = $1`)
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	product, err := repo.GetByID(1)

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

	query := regexp.QuoteMeta(`SELECT id, name, description, price, stock, category_id, created_at, updated_at FROM products WHERE id = $1`)
	mock.ExpectQuery(query).WithArgs(1).WillReturnError(sql.ErrNoRows)

	product, err := repo.GetByID(1)

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

	err = repo.Create(product)

	assert.NoError(t, err)
	assert.Equal(t, 1, product.ID)
	assert.False(t, product.CreatedAt.IsZero())
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

	query := regexp.QuoteMeta(`UPDATE products SET name = $1, price=$2, stock=$3, description=$4, category_id=$5 , updated_at = NOW() WHERE id = $6`)
	mock.ExpectExec(query).
		WithArgs(product.Name, product.Price, product.Stock, product.Description, product.CategoryID, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(1, product)

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

	query := regexp.QuoteMeta(`UPDATE products SET name = $1, price=$2, stock=$3, description=$4, category_id=$5 , updated_at = NOW() WHERE id = $6`)
	mock.ExpectExec(query).
		WithArgs(product.Name, product.Price, product.Stock, product.Description, product.CategoryID, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Update(1, product)

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

	err = repo.Delete(1)

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

	err = repo.Delete(1)

	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
}

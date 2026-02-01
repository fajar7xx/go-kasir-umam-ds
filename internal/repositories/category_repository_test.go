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

func TestCategoryRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	now := time.Now()
	foodDesc := "Food Category"
	bevDesc := "Beverage Category"
	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
		AddRow(1, "Food", &foodDesc, now, now).
		AddRow(2, "Beverage", &bevDesc, now, now)

	query := regexp.QuoteMeta(`SELECT id, name, description, created_at, updated_at FROM categories`)
	mock.ExpectQuery(query).WillReturnRows(rows)

	categories, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Food", categories[0].Name)
	assert.Equal(t, "Beverage", categories[1].Name)
}

func TestCategoryRepository_GetAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	query := regexp.QuoteMeta(`SELECT id, name, description, created_at, updated_at FROM categories`)
	mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)

	categories, err := repo.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, categories)
}

func TestCategoryRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	now := time.Now()
	desc := "Food Category"
	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
		AddRow(1, "Food", &desc, now, now)

	query := regexp.QuoteMeta(`SELECT id, name, description, created_at, updated_at FROM categories where id = $1`)
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	category, err := repo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, 1, category.ID)
	assert.Equal(t, "Food", category.Name)
}

func TestCategoryRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	query := regexp.QuoteMeta(`SELECT id, name, description, created_at, updated_at FROM categories where id = $1`)
	mock.ExpectQuery(query).WithArgs(1).WillReturnError(sql.ErrNoRows)

	category, err := repo.GetByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Equal(t, "Category not found", err.Error())
	assert.Nil(t, category)
}

func TestCategoryRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	now := time.Now()
	desc := "Food Category"
	category := &models.Category{
		Name:        "Food",
		Description: &desc,
	}

	query := regexp.QuoteMeta(`INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, created_at, updated_at`)

	// Mock INSERT RETURNING
	mock.ExpectQuery(query).
		WithArgs(category.Name, category.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, now, now))

	// Mock GetByID call after Create
	getQuery := regexp.QuoteMeta(`SELECT id, name, description, created_at, updated_at FROM categories where id = $1`)
	mock.ExpectQuery(getQuery).WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
			AddRow(1, "Food", &desc, now, now),
	)

	result, err := repo.Create(context.Background(), category)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Food", result.Name)
}

func TestCategoryRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	desc := "Food Desc Updated"
	category := &models.Category{
		Name:        "Food Updated",
		Description: &desc,
	}

	query := regexp.QuoteMeta(`UPDATE categories SET name=$1, description=$2, updated_at=NOW() WHERE id=$3`)
	mock.ExpectExec(query).
		WithArgs(category.Name, category.Description, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), 1, category)

	assert.NoError(t, err)
}

func TestCategoryRepository_Update_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	desc := "Food Desc Updated"
	category := &models.Category{
		Name:        "Food Updated",
		Description: &desc,
	}

	query := regexp.QuoteMeta(`UPDATE categories SET name=$1, description=$2, updated_at=NOW() WHERE id=$3`)
	mock.ExpectExec(query).
		WithArgs(category.Name, category.Description, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Update(context.Background(), 1, category)

	assert.Error(t, err)
	assert.Equal(t, "category not found", err.Error())
}

func TestCategoryRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	query := regexp.QuoteMeta(`DELETE FROM categories where id=$1`)
	mock.ExpectExec(query).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), 1)

	assert.NoError(t, err)
}

func TestCategoryRepository_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	query := regexp.QuoteMeta(`DELETE FROM categories where id=$1`)
	mock.ExpectExec(query).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(context.Background(), 1)

	assert.Error(t, err)
	assert.Equal(t, "category not found", err.Error())
}

func TestCategoryRepository_GetAll_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	now := time.Now()
	desc := "Food Category"

	// Return rows but inject error during iteration
	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
		AddRow(1, "Food", &desc, now, now).
		RowError(0, errors.New("row iteration error"))

	query := regexp.QuoteMeta(`SELECT id, name, description, created_at, updated_at FROM categories`)
	mock.ExpectQuery(query).WillReturnRows(rows)

	categories, err := repo.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, categories)
	assert.Contains(t, err.Error(), "row iteration error")
}

func TestCategoryRepository_Update_RowsAffectedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	desc := "Food Updated"
	category := &models.Category{Name: "Food Updated", Description: &desc}

	query := regexp.QuoteMeta(`UPDATE categories SET name=$1, description=$2, updated_at=NOW() WHERE id=$3`)

	// Return error on RowsAffected()
	mock.ExpectExec(query).
		WithArgs(category.Name, category.Description, 1).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

	err = repo.Update(context.Background(), 1, category)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rows affected error")
}

func TestCategoryRepository_Delete_RowsAffectedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := NewCategoryRepository(db)

	query := regexp.QuoteMeta(`DELETE FROM categories where id=$1`)

	mock.ExpectExec(query).
		WithArgs(1).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

	err = repo.Delete(context.Background(), 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rows affected error")
}

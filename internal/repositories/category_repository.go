package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fajar7xx/go-kasir-umam-ds/models"
)

type CategoryRepositoryInterface interface {
	GetAll(ctx context.Context) ([]models.CategoryResponse, error)
	GetByID(ctx context.Context, id int) (*models.CategoryResponse, error)
	Create(ctx context.Context, category *models.Category) (*models.CategoryResponse, error)
	Update(ctx context.Context, id int, category *models.Category) error
	Delete(ctx context.Context, id int) error
}

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepositoryInterface {
	return &CategoryRepository{
		db: db,
	}
}

func (repo *CategoryRepository) GetAll(ctx context.Context) ([]models.CategoryResponse, error) {
	query := `SELECT
				id, name, description, created_at, updated_at
				FROM categories`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.CategoryResponse, 0, 20)
	for rows.Next() {
		var category models.CategoryResponse
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (repo *CategoryRepository) GetByID(ctx context.Context, id int) (*models.CategoryResponse, error) {
	query := `SELECT
				id, name, description, created_at, updated_at
			FROM categories
			where id = $1`

	var category models.CategoryResponse
	err := repo.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Category not found")
		}
		return nil, err
	}

	return &category, nil
}

func (repo *CategoryRepository) Create(ctx context.Context, category *models.Category) (*models.CategoryResponse, error) {
	query := `INSERT INTO categories
				(name, description)
				VALUES
				($1, $2)
				RETURNING id, created_at, updated_at`

	err := repo.db.QueryRowContext(ctx, query,
		category.Name,
		category.Description,
	).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(ctx, category.ID)
}

func (repo *CategoryRepository) Update(ctx context.Context, id int, category *models.Category) error {
	query := `UPDATE categories
			SET name=$1, description=$2, updated_at=NOW()
			WHERE id=$3`

	result, err := repo.db.ExecContext(ctx, query,
		category.Name,
		category.Description,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category not found")
	}

	return nil
}

func (repo *CategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories where id=$1`
	result, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category not found")
	}

	return nil
}

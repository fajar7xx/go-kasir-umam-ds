package repositories

import (
	"database/sql"
	"errors"
	"fajar7xx/go-kasir-umam-ds/models"
)

type CategoryRepositoryInterface interface {
	GetAll() ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	Create(category *models.Category) error
	Update(id int, category *models.Category) error
	Delete(id int) error
}

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepositoryInterface {
	return &CategoryRepository{
		db: db,
	}
}

func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	query := `SELECT
		id, name, description, created_at, updated_at
		FROM categories`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
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

	return categories, nil
}

func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := `SELECT
				id, name, description, created_at, updated_at
			FROM categories
			where id = $1`

	row := repo.db.QueryRow(query, id)

	var category models.Category
	err := row.Scan(
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

func (repo *CategoryRepository) Create(category *models.Category) error {
	query := `INSERT INTO categories
				(name, description)
				VALUES
				($1, $2)
				RETURNING id, created_at, updated_at`

	err := repo.db.QueryRow(query,
		category.Name,
		category.Description,
	).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *CategoryRepository) Update(id int, category *models.Category) error {
	query := `UPDATE categories
			SET name=$1, description=$2, updated_at=NOW()
			WHERE id=$3`

	result, err := repo.db.Exec(query,
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
		return errors.New("categori not found")
	}

	return nil
}

func (repo *CategoryRepository) Delete(id int) error {
	query := `DELETE FROM categories where id=$1`
	result, err := repo.db.Exec(query, id)
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

package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fajar7xx/go-kasir-umam-ds/models"
)

// 1. ini adalah kontraknya
// siapapun yang ingin menjadi repository product harus punya 5 kemampuan ini.
type ProductRepositoryInterface interface {
	GetAll(ctx context.Context) ([]models.ProductResponse, error)
	GetByID(ctx context.Context, id int) (*models.ProductResponse, error)
	Create(ctx context.Context, product *models.Product) (*models.ProductResponse, error)
	Update(ctx context.Context, id int, product *models.Product) error
	Delete(ctx context.Context, id int) error
}

// 2. ini adalah konkret (si pelakunya)
type ProductRepository struct {
	db *sql.DB
}

// constructor mengembalikan pointer ke struct,
// tapi struc ini secara implisit sudahg memenuhi interface diatas
func NewProductRepository(db *sql.DB) ProductRepositoryInterface {
	return &ProductRepository{
		db: db,
	}
}

func (repo *ProductRepository) GetAll(ctx context.Context) ([]models.ProductResponse, error) {
	query := `select
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
				  join categories c on p.category_id = c.id;`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.ProductResponse, 0, 20) //pre allocate capacity
	for rows.Next() {
		var p models.ProductResponse
		var categoryID int
		var categoryName string
		var categoryDescription *string

		// pastikan urutan scan sesuai dengan urutan select
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.CategoryID,
			&p.CreatedAt,
			&p.UpdatedAt,
			&categoryID,
			&categoryName,
			&categoryDescription)
		if err != nil {
			return nil, err
		}

		p.Category = models.CategorySummary{
			ID:          categoryID,
			Name:        categoryName,
			Description: categoryDescription,
		}

		products = append(products, p)
	}

	// Check error yang terjadi selama iterasi rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (repo *ProductRepository) GetByID(ctx context.Context, id int) (*models.ProductResponse, error) {
	query := `select
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
				where p.id = $1`

	// row := repo.db.QueryRow(query, id)

	var p models.ProductResponse
	var categoryID int
	var categoryName string
	var categoryDescription *string

	// err := row.Scan(
	// 	&p.ID,
	// 	&p.Name,
	// 	&p.Description,
	// 	&p.Price,
	// 	&p.Stock,
	// 	&p.CategoryID,
	// 	&p.CreatedAt,
	// 	&p.UpdatedAt,
	// )

	// QueryRowContext untuk single row + context
	err := repo.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&p.CreatedAt,
		&p.UpdatedAt,
		&categoryID,
		&categoryName,
		&categoryDescription,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// kita bisa return error khusus atau error bawaan sql
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	p.Category = models.CategorySummary{
		ID:          categoryID,
		Name:        categoryName,
		Description: categoryDescription,
	}

	return &p, nil
}

func (repo *ProductRepository) Create(ctx context.Context, product *models.Product) (*models.ProductResponse, error) {
	query := `INSERT INTO products
				(name, price, stock, description, category_id)
			VALUES
				($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at`

	// err := repo.db.QueryRow(query,
	// 	product.Name,
	// 	product.Price,
	// 	product.Stock,
	// 	product.Description,
	// 	product.CategoryID,
	// ).Scan(&product.ID,
	// 	&product.CreatedAt,
	// 	&product.UpdatedAt)

	// QueryRowContext untuk INSERT ... RETURNING
	err := repo.db.QueryRowContext(ctx, query,
		product.Name,
		product.Price,
		product.Stock,
		product.Description,
		product.CategoryID,
	).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return repo.GetByID(ctx, product.ID)
}

func (repo *ProductRepository) Update(ctx context.Context, id int, product *models.Product) error {
	query := `UPDATE products
				SET
				name = $1,
				price=$2,
				stock=$3,
				description=$4,
				category_id=$5,
				updated_at = NOW()
				WHERE id = $6`

	// result, err := repo.db.Exec(query,
	// 	product.Name,
	// 	product.Price,
	// 	product.Stock,
	// 	product.Description,
	// 	product.CategoryID,
	// 	id,
	// )

	// ExecContext untuk UPDATE
	result, err := repo.db.ExecContext(ctx, query,
		product.Name,
		product.Price,
		product.Stock,
		product.Description,
		product.CategoryID,
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
		return errors.New("product not found")
	}

	return nil
}

func (repo *ProductRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE from products where id = $1`

	// result, err := repo.db.Exec(query, id)
	// // ExecContext untuk DELETE
	result, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

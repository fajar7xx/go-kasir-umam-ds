package repositories

import (
	"database/sql"
	"errors"
	"fajar7xx/go-kasir-umam-ds/models"
)

// 1. ini adalah kontraknya
// siapapun yang ingin menjadi repository product harus punya 5 kemampuan ini.
type ProductRepositoryInterface interface {
	GetAll() ([]models.Product, error)
	GetByID(id int) (*models.Product, error)
	Create(product *models.Product) error
	Update(id int, product *models.Product) error
	Delete(id int) error
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

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	query := `SELECT
			id, name, description, price, stock, category_id, created_at, updated_at
			FROM products`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		// pastikan urutan scan sesuai dengan urutan select
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `SELECT
				id, name, description, price, stock, category_id, created_at, updated_at
			FROM products
			WHERE id = $1`

	row := repo.db.QueryRow(query, id)

	var p models.Product
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// kita bisa return error khusus atau error bawaan sql
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := `INSERT INTO products
				(name, price, stock, description, category_id)
			VALUES
				($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at`

	err := repo.db.QueryRow(query,
		product.Name,
		product.Price,
		product.Stock,
		product.Description,
		product.CategoryID,
	).Scan(&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (repo *ProductRepository) Update(id int, product *models.Product) error {
	query := `UPDATE products
				SET name = $1, price=$2, stock=$3, description=$4, category_id=$5 , updated_at = NOW()
				WHERE id = $6`

	result, err := repo.db.Exec(query,
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

func (repo *ProductRepository) Delete(id int) error {
	query := `DELETE from products where id = $1`
	result, err := repo.db.Exec(query, id)
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

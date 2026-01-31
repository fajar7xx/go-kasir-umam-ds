package services

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/internal/repositories"
	"fajar7xx/go-kasir-umam-ds/models"
)

// 1. Definisikan Interface Service (KONTRAK)
// Ini yang akan dipanggil oleh Handler nantinya.
type ProductServiceInterface interface {
	GetAll(ctx context.Context) ([]models.ProductResponse, error)
	GetByID(ctx context.Context, id int) (*models.ProductResponse, error)
	Create(ctx context.Context, product *models.Product) (*models.ProductResponse, error)
	Update(ctx context.Context, id int, product *models.Product) (*models.ProductResponse, error)
	Delete(ctx context.Context, id int) error
}

// 2. Struct Implementasi (Concrete)
type ProductService struct {
	// productRepo *repositories.ProductRepository
	// BEST PRACTICE: Gunakan Interface, bukan struct konkret (*ProductRepository).
	// Ini memungkinkan kita mengganti repo dengan Mock saat Unit Testing.
	productRepo repositories.ProductRepositoryInterface
}

// 3. Constructor
// Perhatikan: Return type-nya sekarang adalah INTERFACE, bukan struct pointer.
//
//	func NewProductService(productRepo *repositories.ProductRepository) *ProductService {
//		return &ProductService{
//			productRepo: productRepo,
//		}
//	}
func NewProductService(productRepo repositories.ProductRepositoryInterface) ProductServiceInterface {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (serv *ProductService) GetAll(ctx context.Context) ([]models.ProductResponse, error) {
	return serv.productRepo.GetAll(ctx)
}

func (serv *ProductService) GetByID(ctx context.Context, id int) (*models.ProductResponse, error) {
	return serv.productRepo.GetByID(ctx, id)
}

func (serv *ProductService) Create(ctx context.Context, product *models.Product) (*models.ProductResponse, error) {
	return serv.productRepo.Create(ctx, product)
}

func (serv *ProductService) Update(ctx context.Context, id int, product *models.Product) (*models.ProductResponse, error) {
	err := serv.productRepo.Update(ctx, id, product)
	if err != nil {
		return nil, err
	}

	updatedProduct, err := serv.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (serv *ProductService) Delete(ctx context.Context, id int) error {
	return serv.productRepo.Delete(ctx, id)
}

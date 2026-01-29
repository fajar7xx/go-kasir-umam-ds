package services

import (
	"fajar7xx/go-kasir-umam-ds/internal/repositories"
	"fajar7xx/go-kasir-umam-ds/models"
)

// 1. Definisikan Interface Service (KONTRAK)
// Ini yang akan dipanggil oleh Handler nantinya.
type ProductServiceInterface interface {
	GetAll() ([]models.Product, error)
	GetByID(id int) (*models.Product, error)
	Create(product *models.Product) error
	Update(id int, product *models.Product) (*models.Product, error)
	Delete(id int) error
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

func (serv *ProductService) GetAll() ([]models.Product, error) {
	return serv.productRepo.GetAll()
}

func (serv *ProductService) GetByID(id int) (*models.Product, error) {
	return serv.productRepo.GetByID(id)
}

func (serv *ProductService) Create(product *models.Product) error {
	return serv.productRepo.Create(product)
}

func (serv *ProductService) Update(id int, product *models.Product) (*models.Product, error) {
	// return serv.productRepo.Update(id, product)
	err := serv.productRepo.Update(id, product)
	if err != nil {
		return nil, err
	}

	updatedProduct, err := serv.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (serv *ProductService) Delete(id int) error {
	return serv.productRepo.Delete(id)
}

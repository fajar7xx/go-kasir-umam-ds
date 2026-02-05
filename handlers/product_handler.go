package handlers

import (
	"context"
	"encoding/json"
	"fajar7xx/go-kasir-umam-ds/internal/services"
	"fajar7xx/go-kasir-umam-ds/models"
	"fajar7xx/go-kasir-umam-ds/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// producthandler mengelola semua endpoint
type ProductHandler struct {
	// nanti bisa ditambah dependency seperti DB, Logger, dll
	productService services.ProductService
}

// newproducthandler membuat instance baru producthandler
func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		utils.SendError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPatch, http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		utils.SendError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Buat context dengan timeout 5 detik
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	// ambil query parameter dari request & sanitize
	queryParams := r.URL.Query()
	name := strings.TrimSpace(queryParams.Get("name"))
	if len(name) > maxNameLength {
		utils.SendError(w, "INVALID_INPUT", "search query is too long", http.StatusBadRequest)
		return
	}

	products, err := h.productService.GetAll(ctx, name)
	if err != nil {
		log.Printf("Get all products error: %v", err)
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT", "Request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "INTERNAL_ERROR", "Failed to get products", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, products, http.StatusOK)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIdFromPath(r, "id")
	if err != nil {
		utils.SendError(w, "INVALID_ID", "invalid product ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	product, err := h.productService.GetByID(ctx, id)
	if err != nil {
		log.Printf("GetByID(%d) error:  %v", id, err)

		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT", "Request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "PRODUCT_NOT_FOUND", "product not found", http.StatusNotFound)
		return
	}

	utils.SendSuccess(w, product, http.StatusOK)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		utils.SendError(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}

	// Server requests: Request body akan otomatis di-close oleh server setelah handler selesai
	// Client requests: Response body WAJIB di-close manual (ini yang sering bikin bingung)
	defer r.Body.Close()

	if err := validateProduct(&newProduct); err != nil {
		utils.SendError(w, "VALIDATION_ERROR", err.Error(), http.StatusBadRequest)
		return
	}

	// Context dengan timeout
	// Kode ini adalah pattern wajib untuk mencegah operasi database/API yang "macet" atau terlalu lama,
	// supaya aplikasi tidak hang.
	// r.Context() = context dari HTTP request (otomatis cancel kalau user cabut koneksi)
	// 5*time.Second = timeout maksimal
	// ctx = context baru yang punya timeout?
	// cancel = fungsi untuk stop paksa (kalau perlu)
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	// Artinya: "Apapun yang terjadi, panggil cancel() saat function ini selesai"
	// Ini wajib untuk cleanup resource (timer, goroutine internal) yang dipakai context
	defer cancel()

	createdProduct, err := h.productService.Create(ctx, &newProduct)
	if err != nil {
		log.Printf("Create product error: %v", err)
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT", "Request Timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "CREATE_FAILED", "Failed to create product", http.StatusBadRequest)
		return
	}

	utils.SendSuccess(w, createdProduct, http.StatusCreated)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIdFromPath(r, "id")
	if err != nil {
		utils.SendError(w, "INVALID_ID", "invalid product ID format", http.StatusBadRequest)
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		utils.SendError(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// simple validation
	if err := validateProduct(&product); err != nil {
		utils.SendError(w, "VALIDATION_ERROR", err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	updatedProduct, err := h.productService.Update(ctx, id, &product)
	if err != nil {
		log.Printf("Update product(%d) error: %v", id, err)

		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT_ERROR", "Request timed out", http.StatusRequestTimeout)
			return
		}
		utils.SendError(w, "UPDATE_FAILED", "failed to update product", http.StatusBadRequest)
		return
	}

	utils.SendSuccess(w, updatedProduct, http.StatusOK)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIdFromPath(r, "id")
	if err != nil {
		utils.SendError(w, "INVALID_ID", "invalid product ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	err = h.productService.Delete(ctx, id)
	if err != nil {
		log.Printf("Delete product(%d) error: %v", id, err)

		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT_ERROR", "Request timed out", http.StatusRequestTimeout)
			return
		}
		utils.SendError(w, "DELETE_FAILED", "failed to delete product", http.StatusBadRequest)
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "product successfully deleted",
	}, http.StatusOK)
}

func validateProduct(product *models.Product) error {
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}

	if len(product.Name) > maxNameLength {
		return fmt.Errorf("product name is too long")
	}

	if product.CategoryID == 0 {
		return fmt.Errorf("category ID is required")
	}

	if product.Price <= 0 {
		return fmt.Errorf("product price must be greater than 0")
	}

	if product.Stock < 0 {
		return fmt.Errorf("product stock must be greater than 0")
	}

	return nil
}

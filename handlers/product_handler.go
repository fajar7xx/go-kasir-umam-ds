package handlers

import (
	"context"
	"encoding/json"
	"fajar7xx/go-kasir-umam-ds/internal/services"
	"fajar7xx/go-kasir-umam-ds/models"
	"fajar7xx/go-kasir-umam-ds/utils"
	"net/http"
	"time"
)

// producthandler mengelola semua endpoint
type ProductHandler struct {
	// nanti bisa ditambah dependency seperti DB, Logger, dll
	productService services.ProductServiceInterface
}

// newproducthandler membuat instance baru producthandler
func NewProductHandler(productService services.ProductServiceInterface) *ProductHandler {
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
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	products, err := h.productService.GetAll(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT", "Request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	product, err := h.productService.GetByID(ctx, id)
	if err != nil {
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

	// simple validation
	if newProduct.Name == "" {
		utils.SendError(w, "VALIDATION_ERROR", "Product name is required", http.StatusBadRequest)
		return
	}

	if newProduct.CategoryID == 0 {
		utils.SendError(w, "VALIDATION_ERROR", "Category ID is required", http.StatusBadRequest)
		return
	}

	if newProduct.Price <= 0 {
		utils.SendError(w, "VALIDATION_ERROR", "Product price must be greater than 0", http.StatusBadRequest)
		return
	}

	if newProduct.Stock <= 0 {
		utils.SendError(w, "VALIDATION_ERROR", "Product stock must be greater than 0", http.StatusBadRequest)
		return
	}

	// Context dengan timeout
	// Kode ini adalah pattern wajib untuk mencegah operasi database/API yang "macet" atau terlalu lama,
	// supaya aplikasi tidak hang.
	// r.Context() = context dari HTTP request (otomatis cancel kalau user cabut koneksi)
	// 5*time.Second = timeout maksimal
	// ctx = context baru yang punya timeout?
	// cancel = fungsi untuk stop paksa (kalau perlu)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	// Artinya: "Apapun yang terjadi, panggil cancel() saat function ini selesai"
	// Ini wajib untuk cleanup resource (timer, goroutine internal) yang dipakai context
	defer cancel()

	createdProduct, err := h.productService.Create(ctx, &newProduct)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT", "Request Timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "CREATE_FAILED", err.Error(), http.StatusBadRequest)
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
	if product.Name == "" {
		utils.SendError(w, "VALIDATION_ERROR", "Product name is required", http.StatusBadRequest)
		return
	}

	if product.CategoryID == 0 {
		utils.SendError(w, "VALIDATION_ERROR", "Category ID is required", http.StatusBadRequest)
		return
	}

	if product.Price <= 0 {
		utils.SendError(w, "VALIDATION_ERROR", "Product price must be greater than 0", http.StatusBadRequest)
		return
	}

	if product.Stock <= 0 {
		utils.SendError(w, "VALIDATION_ERROR", "Product stock must be greater than 0", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	updatedProduct, err := h.productService.Update(ctx, id, &product)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT_ERROR", "Request timed out", http.StatusRequestTimeout)
			return
		}
		utils.SendError(w, "UPDATE_FAILED", err.Error(), http.StatusBadRequest)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = h.productService.Delete(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT_ERROR", "Request timed out", http.StatusRequestTimeout)
			return
		}
		utils.SendError(w, "DELETE_FAILED", err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "product successfully deleted",
	}, http.StatusOK)
}

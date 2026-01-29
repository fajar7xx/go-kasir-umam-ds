package handlers

import (
	"encoding/json"
	"fajar7xx/go-kasir-umam-ds/internal/services"
	"fajar7xx/go-kasir-umam-ds/models"
	"fajar7xx/go-kasir-umam-ds/utils"
	"net/http"
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")

	products, err := h.productService.GetAll()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		utils.SendError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
		return
	}

	// json.NewEncoder(w).Encode(products)
	utils.SendSuccess(w, products, http.StatusOK)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// idStr := r.PathValue("id")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	http.Error(w, "Invalid Request of product id", http.StatusBadRequest)
	// 	return
	// }
	id, err := utils.ParseIdFromPath(r, "id")
	if err != nil {
		utils.SendError(w, "INVALID_ID", "invalid product ID format", http.StatusBadRequest)
		return
	}

	product, err := h.productService.GetByID(id)
	if err != nil {
		// http.Error(w, "invalid product ID", http.StatusNotFound)
		utils.SendError(w, "PRODUCT_NOT_FOUND", "product not found", http.StatusNotFound)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(product)
	utils.SendSuccess(w, product, http.StatusOK)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		// http.Error(w, "Invalid request", http.StatusBadRequest)
		utils.SendError(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.productService.Create(&newProduct)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		utils.SendError(w, "CREATE_FAILED", err.Error(), http.StatusBadRequest)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(newProduct)
	utils.SendSuccess(w, newProduct, http.StatusCreated)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	// idStr := r.PathValue("id")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	http.Error(w, "invalid request of product id", http.StatusBadRequest)
	// 	return
	// }
	id, err := utils.ParseIdFromPath(r, "id")
	if err != nil {
		utils.SendError(w, "INVALID_ID", "invalid product ID format", http.StatusBadRequest)
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		// http.Error(w, "Invalid request", http.StatusBadRequest)
		utils.SendError(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}

	updatedProduct, err := h.productService.Update(id, &product)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		utils.SendError(w, "UPDATE_FAILED", err.Error(), http.StatusBadRequest)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(updatedProduct)
	utils.SendSuccess(w, updatedProduct, http.StatusOK)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// idStr := r.PathValue("id")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	http.Error(w, "invalid request of product id", http.StatusBadRequest)
	// 	return
	// }
	id, err := utils.ParseIdFromPath(r, "id")
	if err != nil {
		utils.SendError(w, "INVALID_ID", "invalid product ID format", http.StatusBadRequest)
		return
	}

	err = h.productService.Delete(id)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		utils.SendError(w, "DELETE_FAILED", err.Error(), http.StatusBadRequest)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"status":  "OK",
	// 	"message": "product has been successfully deleted",
	// })
	utils.SendSuccess(w, map[string]string{
		"message": "product successfully deleted",
	}, http.StatusOK)
}

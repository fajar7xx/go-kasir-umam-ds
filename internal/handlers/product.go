package handlers

import (
	"encoding/json"
	"fajar7xx/go-kasir-umam-ds/internal/models"
	"net/http"
	"strconv"
)

// producthandler mengelola semua endpoint
type ProductHandler struct {
	// nanti bisa ditambah dependency seperti DB, Logger, dll
}

// newproducthandler membuat instance baru producthandler
func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Products)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Request of product id", http.StatusBadRequest)
		return
	}

	for _, product := range models.Products {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// insert data to new products
	newProduct.ID = len(models.Products) + 1
	models.Products = append(models.Products, newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid request of product id", http.StatusBadRequest)
		return
	}

	var updateProduct models.Product
	err = json.NewDecoder(r.Body).Decode(&updateProduct)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i := range models.Products {
		if models.Products[i].ID == id {
			models.Products[i].Name = updateProduct.Name
			models.Products[i].Description = updateProduct.Description
			models.Products[i].Price = updateProduct.Price
			models.Products[i].Stock = updateProduct.Stock

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.Products[i])
			return
		}
	}

	http.Error(w, "product not found", http.StatusNotFound)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid request of product id", http.StatusBadRequest)
		return
	}

	for i, product := range models.Products {
		if product.ID == id {
			// create new slice
			models.Products = append(models.Products[:i], models.Products[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "OK",
				"message": "Product has been successfully deleted",
			})
			return
		}
	}

	http.Error(w, "product not found", http.StatusNotFound)
}

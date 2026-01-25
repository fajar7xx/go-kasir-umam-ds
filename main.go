package main

import (
	"encoding/json"
	"fajar7xx/go-kasir-umam-ds/internal/models"
	"fmt"
	"net/http"
	"strconv"
)

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Request of product id", http.StatusBadRequest)
		return
	}

	for _, product := range models.Products {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
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

func updateProduct(w http.ResponseWriter, r *http.Request) {
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
			// models.Products[i] = updateProduct
			models.Products[i].Name = updateProduct.Name
			models.Products[i].Description = updateProduct.Description
			models.Products[i].Price = updateProduct.Price
			models.Products[i].Stock = updateProduct.Stock

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(models.Products[i])
			return
		}
	}

	http.Error(w, "product not found", http.StatusNotFound)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
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
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "OK",
				"message": "Product has been successfully deleted",
			})
			return
		}
	}

	http.Error(w, "product not found", http.StatusNotFound)
}

func main() {
	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "API Successfull Running on port: 8080",
		})
	})

	// GET /api/v1/products
	// post /api/v1/products
	http.HandleFunc("/api/v1/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProducts(w, r)
		case http.MethodPost:
			createProduct(w, r)
		}
	})

	// get /api/v1/products/{id}
	// put /api/v1/products/{id}
	// delete /api/v1/products/{id}
	http.HandleFunc("/api/v1/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProduct(w, r)
		case http.MethodPut, http.MethodPatch:
			updateProduct(w, r)
		case http.MethodDelete:
			deleteProduct(w, r)
		}
	})

	// get /api/v1/categories
	// post /api/v1/categories
	// get /api/v1/categories/{id}
	// put /api/v1/categories/{id}
	// delete /api/v1/categories/{id}

	fmt.Println("Server started on port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

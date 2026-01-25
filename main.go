package main

import (
	"encoding/json"
	"fajar7xx/go-kasir-umam-ds/internal/handlers"
	"fmt"
	"net/http"
)

func main() {
	// initialize handler
	productHandler := handlers.NewProductHandler()
	categoryHandler := handlers.NewCategoryHandler()

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
			productHandler.GetAll(w, r)
		case http.MethodPost:
			productHandler.Create(w, r)
		default:
			http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		}
	})

	// get /api/v1/products/{id}
	// put /api/v1/products/{id}
	// delete /api/v1/products/{id}
	http.HandleFunc("/api/v1/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetByID(w, r)
		case http.MethodPut, http.MethodPatch:
			productHandler.Update(w, r)
		case http.MethodDelete:
			productHandler.Delete(w, r)
		default:
			http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		}
	})

	// get /api/v1/categories
	// post /api/v1/categories
	http.HandleFunc("/api/v1/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetALL(w, r)
		case http.MethodPost:
			categoryHandler.Create(w, r)
		default:
			http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		}
	})

	// get /api/v1/categories/{id}
	// put /api/v1/categories/{id}
	// delete /api/v1/categories/{id}
	http.HandleFunc("/api/v1/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetByID(w, r)
		case http.MethodPut, http.MethodPatch:
			categoryHandler.Update(w, r)
		case http.MethodDelete:
			categoryHandler.Delete(w, r)
		default:
			http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started on port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

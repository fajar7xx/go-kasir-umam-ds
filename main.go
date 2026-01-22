package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

// variable global yang bisa di akses dimana maa
var products = []Product{
	{ID: 1, Name: "Indomie Goreng", Price: 3500, Stock: 10},
	{ID: 2, Name: "Indomie Ayam", Price: 4000, Stock: 15},
	{ID: 3, Name: "Indomie Telur", Price: 4500, Stock: 20},
	{ID: 4, Name: "Indomie Soto", Price: 5000, Stock: 25},
	{ID: 5, Name: "Indomie Ayam Bakar", Price: 5500, Stock: 30},
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
	http.HandleFunc("/api/v1/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	// POST /api/v1/products

	fmt.Println("Server started on port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

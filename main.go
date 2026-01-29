package main

import (
	"encoding/json"
	"fajar7xx/go-kasir-umam-ds/config"
	"fajar7xx/go-kasir-umam-ds/handlers"
	"fajar7xx/go-kasir-umam-ds/internal/database"
	"fajar7xx/go-kasir-umam-ds/internal/repositories"
	"fajar7xx/go-kasir-umam-ds/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	// 1. load configuration
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := config.Config{
		Port:   viper.GetString("APP_PORT"),
		DBConn: viper.GetString("SUPABASE_DB_CONN"),
	}

	//2. database setup
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("failed to initialize database: ", err)
	}
	defer db.Close()

	// dependency injection
	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	// initialize handler
	// productHandler := handlers.NewProductHandler()
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
	http.HandleFunc("/api/v1/products", productHandler.HandleProducts)

	// get /api/v1/products/{id}
	// put /api/v1/products/{id}
	// delete /api/v1/products/{id}
	http.HandleFunc("/api/v1/products/{id}", productHandler.HandleProductByID)

	// get /api/v1/categories
	// post /api/v1/categories
	http.HandleFunc("/api/v1/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetALL(w, r)
		case http.MethodPost:
			categoryHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running on", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

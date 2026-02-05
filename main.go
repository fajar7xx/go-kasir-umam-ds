package main

import (
	"fajar7xx/go-kasir-umam-ds/config"
	"fajar7xx/go-kasir-umam-ds/handlers"
	"fajar7xx/go-kasir-umam-ds/internal/database"
	"fajar7xx/go-kasir-umam-ds/internal/repositories"
	"fajar7xx/go-kasir-umam-ds/internal/services"
	"fajar7xx/go-kasir-umam-ds/utils"
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

	categoryRepository := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.SendSuccess(w, map[string]string{
			"message": fmt.Sprintf("API Successfull Running on port: %d", config.Port),
		}, http.StatusOK)
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
	http.HandleFunc("/api/v1/categories", categoryHandler.HandleCategories)

	// get /api/v1/categories/{id}
	// put /api/v1/categories/{id}
	// delete /api/v1/categories/{id}
	http.HandleFunc("/api/v1/categories/{id}", categoryHandler.HandleCategoryByID)

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running on", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

package handlers

import (
	"encoding/json"
	"fajar7xx/go-kasir-umam-ds/internal/models"
	"net/http"
	"strconv"
)

// categoryHandler mengelola semua endpoint
type CategoryHandler struct {
	// nanti bisa ditambah dependency seperti DB, logger, dll
}

// newCategoryHandler membuat instance baru CategoryHandler
func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

func (h *CategoryHandler) GetALL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(models.Categories)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invaliud request of category id", http.StatusBadRequest)
	}

	for _, category := range models.Categories {
		if category.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newCategory models.Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	newCategory.ID = len(models.Categories) + 1
	models.Categories = append(models.Categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid request of category id", http.StatusBadRequest)
		return
	}

	var updateCategory models.Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	for i := range models.Categories {
		if models.Categories[i].ID == id {
			models.Categories[i].Name = updateCategory.Name
			models.Categories[i].Description = updateCategory.Description

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.Categories[i])
			return
		}
	}

	http.Error(w, "categor not found", http.StatusNotFound)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid request of category id", http.StatusBadRequest)
		return
	}

	for i, category := range models.Categories {
		if category.ID == id {
			models.Categories = append(models.Categories[:i], models.Categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "OK",
				"message": "Category has been successfully deleted",
			})
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

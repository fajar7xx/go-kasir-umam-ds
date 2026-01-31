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

// categoryHandler mengelola semua endpoint
type CategoryHandler struct {
	// nanti bisa ditambah dependency seperti DB, logger, dll
	categoryService services.CategoryServiceInterface
}

// newCategoryHandler membuat instance baru CategoryHandler
func NewCategoryHandler(categoryService services.CategoryServiceInterface) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		utils.SendError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
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

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	categories, err := h.categoryService.GetAll(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT_ERROR", "request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, categories, http.StatusOK)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIdFromPath(r, "id")
	if err != nil {
		utils.SendError(w, "INVALID_ID", "invalid category ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	category, err := h.categoryService.GetByID(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT_ERROR", "request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "CATEGORY_NOT_FOUND", "category not found", http.StatusNotFound)
		return
	}

	utils.SendSuccess(w, category, http.StatusOK)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newCategory models.Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		utils.SendError(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// simple validation
	if newCategory.Name == "" {
		utils.SendError(w, "INVALID_REQUEST", "category name is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	createdCategory, err := h.categoryService.Create(ctx, &newCategory)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT_ERROR", "request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "CREATE_FAILED", err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccess(w, createdCategory, http.StatusCreated)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIdFromPath(r, "id")
	if err != nil {
		utils.SendError(w, "INVALID_ID", "invalid category ID format", http.StatusBadRequest)
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		utils.SendError(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// simple validation
	if category.Name == "" {
		utils.SendError(w, "INVALID_REQUEST", "category name is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	updatedCategory, err := h.categoryService.Update(ctx, id, &category)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT_ERROR", "request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "UPDATE_FAILED", err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccess(w, updatedCategory, http.StatusOK)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIdFromPath(r, "id")
	if err != nil {
		utils.SendError(w, "INVALID_ID", "invalid category ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.categoryService.Delete(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT_ERROR", "request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "DELETE_FAILED", err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "category successfully deleted",
	}, http.StatusOK)
}

package handlers

import (
	"context"
	"encoding/json"
	"fajar7xx/go-kasir-umam-ds/internal/services"
	"fajar7xx/go-kasir-umam-ds/models"
	"fajar7xx/go-kasir-umam-ds/utils"
	"log"
	"net/http"
)

type TransactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		utils.SendError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var request models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.SendError(w, "INVALID_REQUEST", "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// validation later

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	transaction, err := h.transactionService.Checkout(ctx, request.Items, false)
	if err != nil {
		log.Printf("create transaction error: %v", err)
		if ctx.Err() == context.DeadlineExceeded {
			utils.SendError(w, "TIMEOUT", "Request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "TRANSACTION_FAILED", "failed to create transaction", http.StatusBadRequest)
		return
	}

	utils.SendSuccess(w, transaction, http.StatusCreated)
}

package services

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/internal/repositories"
	"fajar7xx/go-kasir-umam-ds/models"
)

type TransactionService interface {
	Checkout(ctx context.Context, items []models.CheckoutItem, useLock bool) (*models.Transaction, error)
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
}

func NewTrasactionService(transactionRepo repositories.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
	}
}

func (serv *transactionService) Checkout(ctx context.Context, items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return serv.transactionRepo.CreateTransaction(ctx, items)
}

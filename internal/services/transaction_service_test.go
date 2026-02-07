package services

import (
	"context"
	"errors"
	"fajar7xx/go-kasir-umam-ds/internal/mocks"
	"fajar7xx/go-kasir-umam-ds/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransactionService_Checkout_Success(t *testing.T) {
	mockRepo := new(mocks.TransactionRepositoryMock)
	service := NewTrasactionService(mockRepo)

	now := time.Now()
	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
		{ProductID: 2, Quantity: 1},
	}

	expectedTransaction := &models.Transaction{
		ID:          1,
		TotalAmount: 50000.0,
		Details: []models.TransactionDetail{
			{
				ID:            1,
				TransactionID: 1,
				ProductID:     1,
				ProductName:   "Nasi Goreng",
				Price:         20000.0,
				Quantity:      2,
				SubTotal:      40000.0,
				CreatedAt:     now,
			},
		},
		CreatedAt: now,
	}

	mockRepo.On("CreateTransaction", context.Background(), items).Return(expectedTransaction, nil)

	result, err := service.Checkout(context.Background(), items, false)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 50000.0, result.TotalAmount)
	assert.Len(t, result.Details, 1)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_Checkout_UseLockTrue(t *testing.T) {
	mockRepo := new(mocks.TransactionRepositoryMock)
	service := NewTrasactionService(mockRepo)

	now := time.Now()
	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 1},
	}

	expectedTransaction := &models.Transaction{
		ID:          1,
		TotalAmount: 20000.0,
		Details:     []models.TransactionDetail{},
		CreatedAt:   now,
	}

	// Note: Current implementation always calls CreateTransaction regardless of useLock
	// This test documents the current behavior
	// TODO: Fix service to call CreateTransactionOptimal when useLock=true
	mockRepo.On("CreateTransaction", context.Background(), items).Return(expectedTransaction, nil)

	result, err := service.Checkout(context.Background(), items, true)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_Checkout_UseLockFalse(t *testing.T) {
	mockRepo := new(mocks.TransactionRepositoryMock)
	service := NewTrasactionService(mockRepo)

	now := time.Now()
	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 1},
	}

	expectedTransaction := &models.Transaction{
		ID:          1,
		TotalAmount: 20000.0,
		Details:     []models.TransactionDetail{},
		CreatedAt:   now,
	}

	mockRepo.On("CreateTransaction", context.Background(), items).Return(expectedTransaction, nil)

	result, err := service.Checkout(context.Background(), items, false)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_Checkout_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.TransactionRepositoryMock)
	service := NewTrasactionService(mockRepo)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
	}

	mockRepo.On("CreateTransaction", context.Background(), items).Return(nil, errors.New("product not found"))

	result, err := service.Checkout(context.Background(), items, false)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "product not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_Checkout_ContextCancellation(t *testing.T) {
	mockRepo := new(mocks.TransactionRepositoryMock)
	service := NewTrasactionService(mockRepo)

	items := []models.CheckoutItem{
		{ProductID: 1, Quantity: 2},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	mockRepo.On("CreateTransaction", ctx, items).Return(nil, context.Canceled)

	result, err := service.Checkout(ctx, items, false)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_Checkout_EmptyItems(t *testing.T) {
	mockRepo := new(mocks.TransactionRepositoryMock)
	service := NewTrasactionService(mockRepo)

	items := []models.CheckoutItem{}

	// Note: Current implementation doesn't validate empty items at service level
	// Test documents current behavior - validation might be added later
	mockRepo.On("CreateTransaction", context.Background(), items).Return(nil, errors.New("no items to process"))

	result, err := service.Checkout(context.Background(), items, false)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

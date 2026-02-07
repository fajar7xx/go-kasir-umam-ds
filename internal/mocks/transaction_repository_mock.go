package mocks

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/models"

	"github.com/stretchr/testify/mock"
)

type TransactionRepositoryMock struct {
	mock.Mock
}

func (m *TransactionRepositoryMock) CreateTransaction(ctx context.Context, items []models.CheckoutItem) (*models.Transaction, error) {
	args := m.Called(ctx, items)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *TransactionRepositoryMock) CreateTransactionOptimal(ctx context.Context, items []models.CheckoutItem) (*models.Transaction, error) {
	args := m.Called(ctx, items)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

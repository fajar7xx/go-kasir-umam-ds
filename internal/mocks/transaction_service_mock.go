package mocks

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/models"

	"github.com/stretchr/testify/mock"
)

type TransactionServiceMock struct {
	mock.Mock
}

func (m *TransactionServiceMock) Checkout(ctx context.Context, items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	args := m.Called(ctx, items, useLock)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

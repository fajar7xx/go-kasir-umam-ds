package models

import "time"

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount float64             `json:"total_amount"`
	Details     []TransactionDetail `json:"details"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   *time.Time          `json:"updated_at"`
}

type TransactionDetail struct {
	ID            int        `json:"id"`
	TransactionID int        `json:"transaction_id"`
	ProductID     int        `json:"product_id"`
	ProductName   string     `json:"product_name,omitempty"`
	Price         float64    `json:"price"`
	Quantity      int        `json:"quantity"`
	SubTotal      float64    `json:"subtotal"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

package models

import "time"

type Product struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"` //accept null
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	CategoryID  int        `json:"category_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// dummy data
var Products = []Product{
	{ID: 1, Name: "Indomie Goreng", Description: "Indomie Goreng", Price: 3500, Stock: 100},
	{ID: 2, Name: "Indomie Soto", Description: "Indomie Soto", Price: 3600, Stock: 100},
	{ID: 3, Name: "Jus Jeruk", Description: "Jus Jeruk Segar", Price: 5000, Stock: 50},
	{ID: 4, Name: "Pisang Goreng", Description: "Pisang goreng crispy", Price: 2500, Stock: 100},
	{ID: 5, Name: "Risol", Description: "Risol dengan taburan umamy", Price: 2500, Stock: 200},
}

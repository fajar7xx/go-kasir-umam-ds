package models

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var Categories = []Category{
	{ID: 1, Name: "Snack", Description: "Snack"},
	{ID: 2, Name: "Noodle", Description: "Noodle"},
	{ID: 3, Name: "Juice", Description: "Juice"},
}

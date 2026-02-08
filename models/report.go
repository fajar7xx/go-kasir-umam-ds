package models

type ReportResponse struct {
	TotalRevenue      float64     `json:"total_revenue"`
	TotalTransactions int         `json:"total_transactions"`
	BestSeller        *BestSeller `json:"best_seller"`
}

type BestSeller struct {
	Name         string `json:"name"`
	QuantitySold int    `json:"quantity_sold"`
}

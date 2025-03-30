package models

type Status string

const (
	Open   Status = "open"
	Closed Status = "closed"
)

type Order struct {
	ID           string      `json:"order_id"`
	CustomerName string      `json:"customer_name"`
	Status       Status      `json:"status"`
	CreatedAt    string      `json:"created_at"`
	TotalPrice   float64     `json:"total_price"`
	Items        []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID           string  `json:"product_id"`
	Quantity            int     `json:"quantity"`
	PriceDuringTheOrder float64 `json:"price_during_the_order"`
}

type Response struct {
	TotalSales float64 `json:"total-sales"`
}

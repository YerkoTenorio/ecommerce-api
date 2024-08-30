package dto

type OrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	TotalAmount float64             `json:"total_amount"`
	OrderItems  []OrderItemResponse `json:"order_items"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}

type OrderItemResponse struct {
	ProductID uint    `json:"product_id"`
	Quantity  uint    `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

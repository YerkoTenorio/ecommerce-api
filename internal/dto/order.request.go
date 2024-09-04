package dto

type UpdateOrderRequest struct {
	OrderItems []OrderItemRequest `json:"order_items"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
	UnitPrice uint `json:"unit_price"`
}

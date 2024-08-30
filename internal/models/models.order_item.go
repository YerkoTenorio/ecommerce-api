package models

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `json:"product"`
	Quantity  uint    `gorm:"not null" json:"quantity"`
	UnitPrice float64 `gorm:"not null" json:"-"`
	Total     float64 `gorm:"not null" json:"total"`
}

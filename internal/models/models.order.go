package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"user"`
	OrderItem   []OrderItem    `json:"order_item"`
	TotalAmount float64        `gorm:"not null" json:"-"`
	Status      string         `gorm:"size:50;not null;default:'pending'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeleteAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

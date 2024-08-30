package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Stock       int            `gorm:"not null" json:"stock"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeleteAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

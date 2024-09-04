package utils

import (
	"ecommerce-api/internal/dto"
	"ecommerce-api/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil

}

// CheckPasswordHash compara un hash de contraseña con su posible contraseña en texto plano

func CheckPasswordHash(password, hash string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}

// Función auxiliar para crear los detalles de la orden
func CreateOrderDetails(orderItems []models.OrderItem, totalAmount float64) dto.OrderDetails {
	details := dto.OrderDetails{
		OrderItems:  make([]dto.OrderItemResponse, len(orderItems)),
		TotalAmount: totalAmount,
	}

	for i, item := range orderItems {
		details.OrderItems[i] = dto.OrderItemResponse{
			ProductID: item.ProductID,
			Quantity:  int(item.Quantity),
			UnitPrice: item.UnitPrice,
		}
	}

	return details
}

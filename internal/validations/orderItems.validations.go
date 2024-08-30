package validations

import (
	"ecommerce-api/internal/models"
	"ecommerce-api/internal/repositories"
	"ecommerce-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func ValidateOrderItem(orderItem models.OrderItem) *fiber.Map {

	//validar que el id del producto no sea nulo

	if orderItem.ProductID == 0 {
		return &fiber.Map{
			"error": "El id del producto es requerido",
		}
	}

	//validar que el id del producto exista

	if err := utils.DB.Where("id = ?", orderItem.ProductID).First(&models.Product{}).Error; err != nil {
		return &fiber.Map{
			"error": "El producto no existe",
		}
	}

	//validar que la cantidad sea mayor a 0

	if orderItem.Quantity <= 0 {
		return &fiber.Map{
			"error": "La cantidad debe ser mayor a 0",
		}
	}

	//validar que el precio sea mayor a 0
	if orderItem.UnitPrice <= 0 {
		return &fiber.Map{
			"error": "El precio debe ser mayor a 0",
		}
	}

	return nil
}

func CalculateTotalAmount(orderItems []models.OrderItem) (float64, error) {
	var totalAmount float64
	for _, item := range orderItems {
		// asume que el precio unitario se recupera del producto en la base de datos
		product, err := repositories.GetProductById(item.ProductID)
		if err != nil {
			return 0, err
		}
		totalAmount += product.Price * float64(item.Quantity)
	}
	return totalAmount, nil
}

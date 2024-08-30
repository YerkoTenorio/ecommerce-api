package validations

import (
	"ecommerce-api/internal/models"
	"ecommerce-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func ValidateOrder(order models.Order) *fiber.Map {

	//validar que el id del usuario no sea nulo

	if order.UserID == 0 {
		return &fiber.Map{
			"error": "El id del usuario es requerido",
		}
	}

	//validar que el id del usuario exista

	if err := utils.DB.Where("id = ?", order.UserID).First(&models.User{}).Error; err != nil {
		return &fiber.Map{
			"error": "El usuario no existe",
		}
	}

	//validar que la orden tenga al menos un item

	if len(order.OrderItem) == 0 {
		return &fiber.Map{
			"error": "La orden debe tener al menos un item",
		}
	}

	//validar que el total_amount sea mayor a 0

	if order.TotalAmount <= 0 {
		return &fiber.Map{
			"error": "El total_amount debe ser mayor a 0",
		}
	}

	// validar los errores de los items

	for _, item := range order.OrderItem {
		if validateErr := ValidateOrderItem(item); validateErr != nil {
			return validateErr
		}
	}

	return nil
}

package validations

import (
	"ecommerce-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

func ValidateProduct(product models.Product) *fiber.Map {
	if len(product.Name) < 3 || len(product.Name) > 100 {
		return &fiber.Map{
			"error": "El nombre del producto debe tener entre 3 y 100 caracteres",
		}
	}

	if product.Price <= 0 {
		return &fiber.Map{
			"error": "El precio del producto debe ser mayor a 0",
		}
	}

	if product.Stock < 0 {
		return &fiber.Map{
			"error": "El stock del producto debe ser mayor o igual a 0",
		}
	}

	if len(product.Description) > 250 {
		return &fiber.Map{
			"error": "La descripción del producto debe tener menos de 250 caracteres",
		}
	}

	return nil
}

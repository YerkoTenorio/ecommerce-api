package validations

import (
	"ecommerce-api/internal/models"
	"ecommerce-api/internal/repositories"
	"ecommerce-api/pkg/utils"
	"fmt"

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

	// crear un slice para almacenar los Ids de los productos

	productIDs := make([]uint, len(orderItems))
	for i, item := range orderItems {
		productIDs[i] = item.ProductID
	}

	// obtener todos los productos en una sola consulta
	fmt.Println(productIDs)
	products, err := repositories.GetProductsByIds(productIDs)
	if err != nil {
		return 0, err
	}

	fmt.Println(products)

	// crear un mapa para almacenar los precios de los productos

	productPrices := make(map[uint]models.Product)
	for _, product := range products {
		productPrices[product.ID] = product
	}
	//fmt.Println(productPrices)

	// calcular el total de la orden

	for _, item := range orderItems {
		product, exists := productPrices[item.ProductID]
		//fmt.Println(productPrices[item.ProductID])
		if !exists {
			return 0, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("El producto con id %d no existe", item.ProductID))
		}
		totalAmount += float64(item.Quantity) * product.Price

	}

	return totalAmount, nil

}

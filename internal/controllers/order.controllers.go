package controllers

import (
	"ecommerce-api/internal/dto"
	"ecommerce-api/internal/models"
	"ecommerce-api/internal/repositories"
	"ecommerce-api/internal/validations"
	"ecommerce-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	//parsear el body a un modelo de orden

	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Datos invalidos",
		})
	}

	//asignar el unit_price de los items

	for i := range order.OrderItem {
		product, err := repositories.GetProductById(order.OrderItem[i].ProductID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "El producto no existe",
			})
		}
		order.OrderItem[i].UnitPrice = product.Price
	}

	//total_amount es la suma de los totales de los items

	totalAmount, err := validations.CalculateTotalAmount(order.OrderItem)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	order.TotalAmount = totalAmount

	//validar la orden

	if validateErr := validations.ValidateOrder(order); validateErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validateErr)
	}

	//crear la orden

	if err := utils.DB.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al crear la orden",
		})
	}

	//crear los items de la orden

	return c.Status(fiber.StatusCreated).JSON(order)
}

func GetOrders(c *fiber.Ctx) error {

	var orders []models.Order

	// obtener las ordenes y pre-cargar los items de la orden
	if err := utils.DB.Preload("OrderItem").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo obtener las ordenes",
		})
	}

	// preparar la respuesta

	var ordersResponse []dto.OrderResponse
	for _, order := range orders {
		var items []dto.OrderItemResponse
		for _, item := range order.OrderItem {
			items = append(items, dto.OrderItemResponse{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			})
		}

		ordersResponse = append(ordersResponse, dto.OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			TotalAmount: order.TotalAmount,
			OrderItems:  items,
			CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
		})

	}

	return c.Status(fiber.StatusOK).JSON(ordersResponse)

}

func GetOrder(c *fiber.Ctx) error {

	id := c.Params("id")

	var order models.Order

	// obtener la orden y pre-cargar los items de la orden
	if err := utils.DB.Preload("OrderItem").First(&order, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Orden no encontrada",
		})
	}

	// preparar la respuesta

	var items []dto.OrderItemResponse
	for _, item := range order.OrderItem {
		items = append(items, dto.OrderItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})

	}

	orderResponse := dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		OrderItems:  items,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return c.Status(fiber.StatusOK).JSON(orderResponse)

}

func UpdateOrder(c *fiber.Ctx) error {

	return c.SendString("Update Order")
}

func DeleteOrder(c *fiber.Ctx) error {
	return c.SendString("Delete Order")
}

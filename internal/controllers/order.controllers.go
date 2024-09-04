package controllers

import (
	"ecommerce-api/internal/dto"
	"ecommerce-api/internal/models"
	"ecommerce-api/internal/repositories"
	"ecommerce-api/internal/validations"
	"ecommerce-api/pkg/utils"
	"fmt"

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
				Quantity:  int(item.Quantity),
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
			Quantity:  int(item.Quantity),
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

	id := c.Params("id")
	var req dto.UpdateOrderRequest

	// Parsear el body a la estructura
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Datos invalidos",
		})
	}

	// Buscar la orden en la base de datos
	var order models.Order
	if err := utils.DB.Preload("OrderItem").First(&order, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Orden no encontrada",
		})
	}

	// Capturar los datos actuales de la orden antes de actualizar
	oldOrderDetails := utils.CreateOrderDetails(order.OrderItem, order.TotalAmount)

	// Iniciar una transacción para asegurar atomicidad
	tx := utils.DB.Begin()

	// Eliminar los ítems existentes de la orden
	if err := tx.Where("order_id = ?", order.ID).Delete(&models.OrderItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al eliminar los ítems de la orden",
		})
	}

	// Obtener los IDs de los productos del request
	productIDs := make([]uint, len(req.OrderItems))
	for i, item := range req.OrderItems {
		productIDs[i] = item.ProductID
	}

	// Recuperar todos los productos necesarios en una sola consulta
	var products []models.Product
	if err := tx.Where("id IN (?)", productIDs).Find(&products).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al recuperar productos",
		})
	}

	// Crear un mapa para acceder a los productos por su ID
	productMap := make(map[uint]models.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	// Actualizar los ítems de la orden con precios actualizados
	updatedItems := make([]models.OrderItem, len(req.OrderItems))
	for i, item := range req.OrderItems {
		product, exists := productMap[item.ProductID]
		if !exists {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("El producto con ID %d no existe", item.ProductID),
			})
		}
		updatedItems[i] = models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  uint(item.Quantity),
			UnitPrice: product.Price,
		}
	}

	// Calcular el monto total utilizando la función optimizada
	totalAmount, err := validations.CalculateTotalAmount(updatedItems)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al calcular el total de la orden",
		})
	}

	// Actualizar la orden con los nuevos ítems y el total calculado
	order.OrderItem = updatedItems
	order.TotalAmount = totalAmount

	// Guardar los cambios en la base de datos
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al guardar la orden",
		})
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al confirmar la transacción",
		})
	}

	// Crear los datos actuales de la orden después de actualizar
	newOrderDetails := utils.CreateOrderDetails(order.OrderItem, order.TotalAmount)

	// Construir la respuesta con los datos antiguos y nuevos
	response := dto.OrderUpdateResponse{
		OldOrder: oldOrderDetails,
		NewOrder: newOrderDetails,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteOrder(c *fiber.Ctx) error {

	id := c.Params("id")
	var order models.Order

	// buscar la orden en la base de datos
	if err := utils.DB.First(&order, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Orden no encontrada",
		})
	}

	//iniciar una transacción para asegurar atomicidad
	tx := utils.DB.Begin()

	// buscar la orden en la base de datos

	if err := tx.Preload("OrderItem").First(&order, id).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Orden no encontrada",
		})
	}

	// eliminar los items de la orden
	if err := tx.Where("order_id = ?", order.ID).Delete(&models.OrderItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al eliminar los items de la orden",
		})
	}

	// eliminar la orden

	if err := tx.Delete(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al eliminar la orden",
		})
	}

	// confirmar la transacción

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al confirmar la transacción",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Orden eliminada exitosamente",
	})
}

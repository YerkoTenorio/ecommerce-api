package controllers

import (
	"ecommerce-api/internal/models"
	"ecommerce-api/internal/validations"
	"ecommerce-api/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {

	var product models.Product
	fmt.Println("recibiendo los datos del producto...")
	if err := c.BodyParser(&product); err != nil {
		fmt.Println("error al parsear los datos del producto", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Datos invalidos",
		})
	}

	// validar los datos del producto

	fmt.Println("datos del producto recibidos despues de parsear", product)

	if validationsErrors := validations.ValidateProduct(product); validationsErrors != nil {
		fmt.Println("error al validar los datos del producto")
		return c.Status(fiber.StatusBadRequest).JSON(validationsErrors)
	}

	// guardar el producto en la base de datos

	if err := utils.DB.Create(&product).Error; err != nil {
		fmt.Println("error al guardar el producto en la base de datos", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo guardar el producto",
		})
	}

	fmt.Println("producto guardado exitosamente")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Producto creado exitosamente",
	})
}

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	if err := utils.DB.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo obtener los productos",
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	if err := utils.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Producto no encontrado",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)

}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	// buscar el producto en la base de datos
	if err := utils.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Producto no encontrado",
		})
	}

	// parsear los datos del producto

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Datos invalidos",
		})
	}

	// validar los datos del producto

	if validationsErrors := validations.ValidateProduct(product); validationsErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationsErrors)
	}

	// guardar los cambios en la base de datos

	if err := utils.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo actualizar el producto",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)

}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	// buscar el producto en la base de datos
	if err := utils.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Producto no encontrado",
		})
	}

	// eliminar el producto de la base de datos

	if err := utils.DB.Delete(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo eliminar el producto",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Producto eliminado exitosamente",
	})

}

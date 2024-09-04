package controllers

import (
	"ecommerce-api/internal/dto"
	"ecommerce-api/internal/models"
	"ecommerce-api/internal/validations"
	"ecommerce-api/pkg/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {

	var req dto.RegisterRequest

	// parsear el body a la estructura

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo parsear el body",
		})
	}

	//validar la solicitud de registro

	if validationErrors := validations.ValidateRegisterRequest(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	// hashear la contrasena

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo hashear la contraseña",
		})
	}

	// crear el usuario y guardarlo en la bd

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", //rol por defecto
	}

	if err := utils.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo crear el usuario",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Usuario creado exitosamente",
	})
}

func Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	// parsear el body a la estructura

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo parsear el body",
		})
	}

	// validar la solicitud de login

	if validationErrors := validations.ValidateLoginRequest(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	// buscar el usuario en la bd

	var user models.User

	if err := utils.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Correo o contraseña incorrectos",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al buscar el usuario",
		})

	}

	// comparar la contrasena

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Correo o contraseña incorrectos",
		})
	}

	// generar el token jwt

	token, err := utils.GenerateJWT(user.ID, os.Getenv("JWT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al generar el token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})

}

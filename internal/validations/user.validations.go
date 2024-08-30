package validations

import (
	"ecommerce-api/internal/dto"
	"ecommerce-api/internal/models"
	"ecommerce-api/pkg/utils"
	"ecommerce-api/pkg/validations"

	"github.com/gofiber/fiber/v2"
)

func ValidateRegisterRequest(req dto.RegisterRequest) *fiber.Map {
	//validar el largo del nombre

	if len(req.Name) < 3 || len(req.Name) > 50 {
		return &fiber.Map{
			"error": "El nombre debe tener entre 3 y 50 caracteres",
		}
	}

	//validar email

	if !validations.IsValidEmail(req.Email) {
		return &fiber.Map{
			"error": "El email no es válido",
		}
	}

	//validar si existe el email en la bd

	var existingUser models.User

	if err := utils.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return &fiber.Map{
			"error": "El email ya está en uso",
		}
	}

	// validar largo de la contrasena

	if len(req.Password) < 6 || len(req.Password) > 50 {
		return &fiber.Map{
			"error": "La contraseña debe tener entre 6 y 50 caracteres",
		}
	}

	// validar si las contrasenas coinciden

	if req.Password != req.ConfirmPassword {
		return &fiber.Map{
			"error": "Las contraseñas no coinciden",
		}
	}

	return nil

}

func ValidateLoginRequest(req dto.LoginRequest) *fiber.Map {
	//validar email

	if req.Email == "" {
		return &fiber.Map{
			"error": "El email es requerido",
		}
	}

	// validar contraseña

	if req.Password == "" {
		return &fiber.Map{
			"error": "La contraseña es requerida",
		}
	}
	return nil

}

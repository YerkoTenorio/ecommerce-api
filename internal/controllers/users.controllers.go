package controllers

import "github.com/gofiber/fiber/v2"

func GetAllUsers(c *fiber.Ctx) error {
	return c.SendString("GetAllUsers")
}

func GetUser(c *fiber.Ctx) error {
	return c.SendString("GetUser")
}

func UpdateUser(c *fiber.Ctx) error {
	return c.SendString("UpdateUser")
}

func DeleteUser(c *fiber.Ctx) error {
	return c.SendString("DeleteUser")
}

func GetProfile(c *fiber.Ctx) error {
	return c.SendString("GetProfile")
}

func UpdateProfile(c *fiber.Ctx) error {
	return c.SendString("UpdateProfile")
}

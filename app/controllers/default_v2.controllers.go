package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/godgodwinter/go-fiber-1-keuangan/app/models"
	"github.com/godgodwinter/go-fiber-1-keuangan/app/services"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func V2_GetAll(c *fiber.Ctx) error {
	service := services.DefaultService{}
	return service.V2_GetAll(c)
}
func V2_GetUserWhereId(c *fiber.Ctx) error {
	service := services.DefaultService{}
	userId := c.Params("userId")
	return service.V2_GetUserWhereId(c, userId)
}

//	func V2_CreateUser(c *fiber.Ctx) error {
//		service := services.DefaultService{}
//		return service.V2_CreateUser(c)
//	}
func V2_CreateUser(c *fiber.Ctx) error {
	var user models.User

	// Parse request body into user struct
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate the user struct using validator library
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation error",
			"errors":  validationErr.Error(),
		})
	}

	service := services.DefaultService{}
	return service.V2_CreateUser(c, user)
}
func V2_UpdateUser(c *fiber.Ctx) error {
	service := services.DefaultService{}
	userId := c.Params("userId")
	var user models.User

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    err.Error(),
			"error":   err.Error(),
		})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    validationErr.Error(),
			"error":   validationErr.Error(),
		})
	}

	return service.V2_UpdateUser(c, userId, user)
}
func V2_DeleteUser(c *fiber.Ctx) error {
	service := services.DefaultService{}
	userId := c.Params("userId")
	return service.V2_DeleteUser(c, userId)
}

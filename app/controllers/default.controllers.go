package controllers

import (
	"github.com/godgodwinter/go-fiber-1-keuangan/app/services"
	"github.com/gofiber/fiber/v2"
)

func GetAll(c *fiber.Ctx) error {
	service := services.DefaultService{}
	return service.GetAll(c)
}
func GetAUser(c *fiber.Ctx) error {
	service := services.DefaultService{}
	return service.GetAUser(c)
}
func CreateUser(c *fiber.Ctx) error {
	service := services.DefaultService{}
	return service.CreateUser(c)
}
func GetAllUser(c *fiber.Ctx) error {
	service := services.DefaultService{}
	return service.GetAllUser(c)
}

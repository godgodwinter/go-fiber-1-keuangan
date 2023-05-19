package controllers

import (
	"github.com/godgodwinter/go-fiber-1-keuangan/app/services"
	"github.com/gofiber/fiber/v2"
)

func GetBasic(c *fiber.Ctx) error {
	service := services.DefaultService{}
	return service.GetBasic(c)
}
func GetAll(c *fiber.Ctx) error {
	service := services.DefaultService{}
	return service.GetAll(c)
}
func GetUserWhereId(c *fiber.Ctx) error {
	service := services.DefaultService{}
	return service.GetUserWhereId(c)
}
func CreateUser(c *fiber.Ctx) error {
	service := services.DefaultService{}
	return service.CreateUser(c)
}
func UpdateUser(c *fiber.Ctx) error {
	service := services.DefaultService{}
	return service.UpdateUser(c)
}
func DeleteUser(c *fiber.Ctx) error {
	service := services.DefaultService{}
	return service.DeleteUser(c)
}

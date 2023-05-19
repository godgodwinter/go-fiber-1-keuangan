package routes

import (
	"fmt"

	"github.com/godgodwinter/go-fiber-1-keuangan/app/configs"
	"github.com/gofiber/fiber/v2"
)

type ExampleStruct struct {
	Name string
	Desc string
}

func DefaultRoutes(app *fiber.App) {

	app_version := configs.EnvAppVersion()
	api_version := app.Group("/api/" + app_version)
	app_routes := api_version.Group("/default")
	app_routes.Get("/index", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => 💸 From: LAX, To: SFO
	})
	app_routes.Get("/nama/:name", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => 💸 From: LAX, To: SFO
	})

	app_routes.Get("/json/:name", func(c *fiber.Ctx) error {
		data := ExampleStruct{
			Name: c.Params("name"),
			Desc: "Ini desc",
		}
		return c.JSON(data)
	})
}

package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/godgodwinter/go-fiber-1-keuangan/app/configs"
	"github.com/godgodwinter/go-fiber-1-keuangan/app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor" //digunakan untuk menangani panic pada saat runtime. Panic terjadi ketika terjadi kesalahan yang tidak dapat ditangani pada saat aplikasi berjalan, dan biasanya menyebabkan program berhenti secara tiba-tiba. go recover digunakan untuk menangani panic sehingga aplikasi tidak berhenti dengan sendirinya.
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	strVar := configs.EnvPort()
	app_port, err := strconv.Atoi(strVar)
	if err != nil {
		// handle error
		log.Fatal("Error env port tidak ditemukan")
	}
	//run database
	configs.ConnectDB()

	// !Custom config
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		Network:       "tcp4",
		ServerHeader:  "Fiber",
		AppName:       "go-fiber-1-keuangan v1.0.1",
	})

	app.Get("/about", func(c *fiber.Ctx) error {
		return c.SendString("about")
	})

	type SomeStruct struct {
		Name string
		Port int
	}

	app.Get("/", func(c *fiber.Ctx) error {
		// Create data struct:
		data := SomeStruct{
			Name: fmt.Sprintf("go-fiber-1-keuangan port: %d", app_port),
			Port: app_port,
		}

		return c.JSON(data)
	})

	app.Use(recover.New())
	routes.DefaultRoutes(app)

	//   !monitor
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	// Menentukan alamat dan port
	addr := fmt.Sprintf("0.0.0.0:%d", app_port)

	// Mulai server HTTP
	app.Listen(addr)

}

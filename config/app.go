package config

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// 🟩 Tambahkan CORS agar Swagger UI bisa fetch API tanpa error
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // boleh ubah jadi "http://localhost:3000" kalau mau lebih aman
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 🟦 Logging middleware tetap dipertahankan
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// 🟨 Kalau kamu ingin serve folder uploads secara global:
	// app.Static("/uploads", "./uploads")

	return app
}

func StartServer(app *fiber.App, port string) {
	log.Fatal(app.Listen(":" + port))
}

package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Inisialisasi Fiber app
	app := fiber.New()

	// Endpoint root "/"
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! Praktikum Go Fiber berjalan.")
	})

	// Jalankan server di port 3000
	log.Fatal(app.Listen(":3000"))
}
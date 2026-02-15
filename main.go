package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	fmt.Println("listening on http://localhost:3000/")

	app.Listen(":3000")

}

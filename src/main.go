package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/utils"
)

func main() {
	app := fiber.New()

	utils.InitDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}

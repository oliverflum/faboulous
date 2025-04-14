package main

import (
	"github.com/gofiber/fiber/v2/log"

	"github.com/oliverflum/faboulous/backend/app"
)

func main() {
	app := app.SetupApp()
	log.Fatal(app.Listen(":3000"))
}

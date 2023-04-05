package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/handler"
	"github.com/oliverflum/faboulous/util"
)

func main() {
	app := fiber.New()

	util.InitDB()

	api := app.Group("/api")

	configApi := api.Group("/config")
	configApi.Get("/feature", handler.ListFeatures)
	configApi.Post("/feature", handler.AddFeature)
	configApi.Get("/feature/:id", handler.GetFeature)
	configApi.Delete("/feature/:id", handler.DeleteFeature)

	log.Fatal(app.Listen(":3000"))
}

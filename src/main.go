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

	featureApi := configApi.Group("/feature")
	featureApi.Get("/", handler.ListFeatures)
	featureApi.Post("/", handler.AddFeature)
	featureApi.Put("/:id", handler.UpdateFeature)
	featureApi.Get("/:id", handler.GetFeature)
	featureApi.Delete("/:id", handler.DeleteFeature)

	// testApi := configApi.Group("/test")
	// testApi.Get("/", handler.ListFeatures)
	// testApi.Post("/", handler.AddFeature)
	// testApi.Put("/", handler.AddFeature)
	// testApi.Get("/:id", handler.GetFeature)
	// testApi.Delete("/:id", handler.DeleteFeature)

	log.Fatal(app.Listen(":3000"))
}

package main

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/db"
	"github.com/oliverflum/faboulous/handler"
)

func main() {
	app := fiber.New()
	log.SetLevel(log.LevelTrace)
	app.Use(logger.New())

	db.InitDB()

	api := app.Group("/api")

	configApi := api.Group("/config")

	featureApi := configApi.Group("/feature")
	featureApi.Get("/", handler.ListFeatures)
	featureApi.Post("/", handler.AddFeature)
	featureApi.Put("/:id", handler.UpdateFeature)
	featureApi.Get("/:id", handler.GetFeature)
	featureApi.Delete("/:id", handler.DeleteFeature)

	testApi := configApi.Group("/test")
	testApi.Get("/", handler.ListTests)
	testApi.Post("/", handler.AddTest)
	testApi.Put("/:id", handler.UpdateTest)
	testApi.Get("/:id", handler.GetTest)
	testApi.Delete("/:id", handler.DeleteTest)

	variantApi := testApi.Group(":testId/variant")
	variantApi.Post("/", handler.AddVariant)
	variantApi.Put("/:id", handler.UpdateVariant)
	variantApi.Delete("/:id", handler.DeleteVariant)

	log.Fatal(app.Listen(":3000"))
}

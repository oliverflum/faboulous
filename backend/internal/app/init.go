package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/oliverflum/faboulous/backend/handler"
	"github.com/oliverflum/faboulous/backend/internal/db"
	"github.com/oliverflum/faboulous/backend/internal/util"
)

func SetupApp() *fiber.App {
	app := fiber.New()

	log.SetLevel(util.GetLogLevel())
	app.Use(logger.New())

	db.InitDB()

	api := app.Group("/api")

	configApi := api.Group("/admin")

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

	variantFeatureApi := variantApi.Group(":variantId/feature")
	variantFeatureApi.Post("/", handler.AddVariantFeature)
	variantFeatureApi.Put("/:id", handler.UpdateVariantFeature)
	variantFeatureApi.Delete("/:id", handler.DeleteVariantFeature)

	return app
}

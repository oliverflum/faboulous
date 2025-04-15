package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/handler"
	"github.com/oliverflum/faboulous/backend/util"
)

func SetupApp() *fiber.App {
	app := fiber.New()

	log.SetLevel(util.GetLogLevel())
	app.Use(logger.New())

	db.InitDB()

	api := app.Group("/api")

	configApi := api.Group("/config")
	configApi.Get("/", handler.GetConfig)
	adminApi := api.Group("/admin")

	// Add the publish route
	adminApi.Post("/publish", handler.Publish)

	featureApi := adminApi.Group("/feature")
	featureApi.Get("/", handler.ListFeatures)
	featureApi.Post("/", handler.AddFeature)
	featureApi.Put("/:featureId", handler.UpdateFeature)
	featureApi.Get("/:featureId", handler.GetFeature)
	featureApi.Delete("/:featureId", handler.DeleteFeature)

	testApi := adminApi.Group("/test")
	testApi.Get("/", handler.ListTests)
	testApi.Post("/", handler.AddTest)
	testApi.Put("/:testId", handler.UpdateTest)
	testApi.Get("/:testId", handler.GetTest)
	testApi.Delete("/:testId", handler.DeleteTest)

	variantApi := testApi.Group(":testId/variant")
	variantApi.Get("/", handler.ListVariants)
	variantApi.Post("/", handler.AddVariant)
	variantApi.Put("/:variantId", handler.UpdateVariant)
	variantApi.Delete("/:variantId", handler.DeleteVariant)

	variantFeatureApi := variantApi.Group(":variantId/variant_feature")
	variantFeatureApi.Post("/", handler.AddVariantFeature)
	variantFeatureApi.Put("/:variantFeatureId", handler.UpdateVariantFeature)
	variantFeatureApi.Delete("/:variantFeatureId", handler.DeleteVariantFeature)

	return app
}

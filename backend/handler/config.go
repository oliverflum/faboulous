package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/service"
)

func GetConfig(c *fiber.Ctx) error {
	testsConfigs, err := service.GetFeatureSet(c.Query("differentiator"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(testsConfigs)
}

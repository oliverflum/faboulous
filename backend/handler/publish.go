package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/service"
)

func Publish(c *fiber.Ctx) error {
	err := service.PublishConfig()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusNoContent).SendString("Cache updated")
}

package handler

import (
	"github.com/gofiber/fiber/v2"
)

func GetConfig(c *fiber.Ctx) error {
	// d := c.Query("d")
	// testsConfigs := db.GetTestConfigs()
	return c.Status(fiber.StatusOK).JSON("FABOULOUS")
}

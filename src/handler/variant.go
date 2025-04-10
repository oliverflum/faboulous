package handler

import (
	"github.com/gofiber/fiber/v2"
)

func ListVariants(c *fiber.Ctx) error {
	// Logic to list all variants
	return c.JSON(fiber.Map{"message": "List of variants"})
}

func AddVariant(c *fiber.Ctx) error {
	// Logic to add a new variant
	return c.JSON(fiber.Map{"message": "Variant added"})
}

func UpdateVariant(c *fiber.Ctx) error {
	// Logic to update an existing variant
	return c.JSON(fiber.Map{"message": "Variant updated"})
}

func DeleteVariant(c *fiber.Ctx) error {
	// Logic to delete a variant
	return c.JSON(fiber.Map{"message": "Variant deleted"})
}
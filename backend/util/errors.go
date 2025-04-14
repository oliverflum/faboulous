package util

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HandleGormError(result *gorm.DB) *fiber.Error {
	if result.Error == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "Record not found")
	}
	return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
}

package util

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HandleGormError(err error) error {
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "Record not found")
	}
	return err
}

func SendErrorRes(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *fiber.Error:
		return c.Status(e.Code).SendString(e.Message)
	default:
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
}

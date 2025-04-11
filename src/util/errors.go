package util

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FabolousError struct {
	Code    int
	Message string
}

func (e FabolousError) Error() string {
	return e.Message
}

func HandleGormError(err error) error {
	if err == gorm.ErrRecordNotFound {
		return FabolousError{
			Code:    fiber.StatusNotFound,
			Message: "Record not found",
		}
	}
	return err
}

func SendErrorRes(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case FabolousError:
		return c.Status(e.Code).SendString(e.Message)
	default:
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
}

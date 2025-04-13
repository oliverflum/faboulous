package util

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(item any) []*ErrorResponse {
	var validate = validator.New()
	var errors []*ErrorResponse
	err := validate.Struct(item)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ParseAndValidatePayload[T any](c *fiber.Ctx, payload *T) *fiber.Error {
	if err := c.BodyParser(payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate the payload
	errors := ValidateStruct(payload)

	if len(errors) > 0 {
		errorsString := ""
		for _, error := range errors {
			errorsString += error.FailedField + " " + error.Tag + " " + error.Value + "\n"
		}
		return fiber.NewError(fiber.StatusBadRequest, errorsString)
	}
	return nil
}

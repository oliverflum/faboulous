package utils

import {
	"github.com/gofiber/fiber/v2"
}

const validate = validator.New()

func ValidateStruct(item interface{}) []*ErrorResponse {
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
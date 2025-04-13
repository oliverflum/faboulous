package util

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	BOOL   = "BOOL"
	STRING = "STRING"
	INT    = "INT"
	FLOAT  = "FLOAT"
)

func GetValueTypeAndString(jsonValue any) (string, string, error) {
	switch v := jsonValue.(type) {
	case bool:
		return BOOL, fmt.Sprintf("%t", v), nil
	case string:
		return STRING, v, nil
	case int, int8, int16, int32, int64:
		return INT, fmt.Sprintf("%d", v), nil
	case float32, float64:
		return FLOAT, fmt.Sprintf("%f", v), nil
	default:
		return "", "", fiber.NewError(fiber.StatusBadRequest, "unsupported value type: "+fmt.Sprintf("%T", jsonValue))
	}
}

func GetJsonValue(value string, valueType string) (any, error) {
	conversionError := fiber.NewError(fiber.StatusInternalServerError, "Invalid value for "+valueType+": "+value)
	switch valueType {
	case BOOL:
		var boolValue bool
		if _, err := fmt.Sscanf(value, "%t", &boolValue); err != nil {
			return nil, conversionError
		}
		return boolValue, nil
	case STRING:
		return value, nil
	case INT:
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err != nil {
			return nil, conversionError
		}
		return intValue, nil
	case FLOAT:
		var floatValue float64
		if _, err := fmt.Sscanf(value, "%f", &floatValue); err != nil {
			return nil, conversionError
		}
		return floatValue, nil
	default:
		return nil, fiber.NewError(fiber.StatusInternalServerError, "unsupported value type: "+valueType)
	}
}

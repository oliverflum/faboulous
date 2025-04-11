package model

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/util"
	"gorm.io/gorm"
)

const (
	BOOL   = "BOOL"
	STRING = "STRING"
	INT    = "INT"
	FLOAT  = "FLOAT"
)

type Feature struct {
	gorm.Model
	Name         string `gorm:"not null;unique"`
	Type         string `gorm:"not null"`
	DefaultValue string `gorm:"not null"`
}

func (feature *Feature) UpdateFromPayload(payload FeaturePayload) error {
	feature.Name = payload.Name
	valueType, defaultValue, err := GetEntityValueAndType(payload.Value)
	if err != nil {
		return err
	}
	feature.Type = valueType
	feature.DefaultValue = defaultValue
	return nil
}

func GetEntityValueAndType(value any) (string, string, error) {
	switch v := value.(type) {
	case bool:
		return BOOL, fmt.Sprintf("%t", v), nil
	case string:
		return STRING, v, nil
	case int, int8, int16, int32, int64:
		return INT, fmt.Sprintf("%d", v), nil
	case float32, float64:
		return FLOAT, fmt.Sprintf("%f", v), nil
	default:
		return "", "", util.FabolousError{
			Code:    fiber.StatusBadRequest,
			Message: "unsupported value type: " + fmt.Sprintf("%T", value),
		}
	}
}

func NewFeatureEntity(feature FeaturePayload) (Feature, error) {
	valueType, defaultValue, err := GetEntityValueAndType(feature.Value)
	if err != nil {
		return Feature{}, errors.New("Could not instantiate feature entity: " + err.Error())
	}
	return Feature{
		Name:         feature.Name,
		Type:         valueType,
		DefaultValue: defaultValue,
	}, nil
}

type FeaturePayload struct {
	Id    uint   `json:"id"`
	Name  string `json:"name" validate:"required"`
	Value any    `json:"value" validate:"required"`
}

func GetPayloadValue(value string, valueType string) (any, error) {
	switch valueType {
	case BOOL:
		return value == "true", nil
	case STRING:
		return value, nil
	case INT:
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err != nil {
			return nil, err
		}
		return intValue, nil
	case FLOAT:
		var floatValue float64
		if _, err := fmt.Sscanf(value, "%f", &floatValue); err != nil {
			return nil, err
		}
		return floatValue, nil
	default:
		return nil, util.FabolousError{
			Code:    fiber.StatusInternalServerError,
			Message: "unsupported value type: " + valueType,
		}
	}
}

func NewFeaturePayload(feature Feature) (FeaturePayload, error) {
	value, err := GetPayloadValue(feature.DefaultValue, feature.Type)

	if err != nil {
		return FeaturePayload{}, errors.New("Could not instantiate feature payload: " + err.Error())
	}

	return FeaturePayload{
		Id:    feature.ID,
		Name:  feature.Name,
		Value: value,
	}, nil
}

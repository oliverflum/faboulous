package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const (
	BOOL   = "BOOL"
	STRING = "STRING"
	INT    = "INT"
	FLOAT  = "FLOAT"
)

type FeatureEntity struct {
	gorm.Model
	Name         string `gorm:"not null;unique"`
	Type         string `gorm:"not null"`
	DefaultValue string `gorm:"not null"`
}

func (feature *FeatureEntity) TableName() string {
	return "features"
}

func getEntityValueAndType(value any) (string, string, error) {
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
		return "", "", fmt.Errorf("unsupported value type: %T", value)
	}
}

func NewFeatureEntity(feature FeaturePayload) (FeatureEntity, error) {
	valueType, defaultValue, err := getEntityValueAndType(feature.Value)
	if err != nil {
		return FeatureEntity{}, errors.New("Could not instantiate feature entity: " + err.Error())
	}
	return FeatureEntity{
		Name:         feature.Name,
		Type:         valueType,
		DefaultValue: defaultValue,
	}, nil
}

type FeaturePayload struct {
	Name  string `json:"name" validate:"required"`
	Value any    `json:"value" validate:"required"`
}

func getPayloadValue(value string, valueType string) (any, error) {
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
		return nil, fmt.Errorf("unsupported value type: %s", valueType)
	}
}

func NewFeaturePayload(feature FeatureEntity) (FeaturePayload, error) {
	value, err := getPayloadValue(feature.DefaultValue, feature.Type)

	if err != nil {
		return FeaturePayload{}, errors.New("Could not instantiate feature payload: " + err.Error())
	}

	return FeaturePayload{
		Name:  feature.Name,
		Value: value,
	}, nil
}

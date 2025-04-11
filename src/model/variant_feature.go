package model

import (
	"fmt"

	"gorm.io/gorm"
)

type VariantFeatureEntity struct {
	gorm.Model
	FeatureID uint
	Type      string
	Value     string
}

func (variantFeature *VariantFeatureEntity) TableName() string {
	return "variant_features"
}

func (variantFeature *VariantFeatureEntity) SetValue(value any) error {
	valueType, stringValue, err := getEntityValueAndType(value)
	if err != nil {
		return fmt.Errorf("could not set variant feature value: %w", err)
	}
	variantFeature.Type = valueType
	variantFeature.Value = stringValue
	return nil
}

func (variantFeature *VariantFeatureEntity) GetValue() (any, error) {
	value, err := getPayloadValue(variantFeature.Value, variantFeature.Type)
	if err != nil {
		return nil, fmt.Errorf("could not get variant feature value: %w", err)
	}
	return value, nil
}

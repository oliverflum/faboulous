package model

import (
	"fmt"

	"gorm.io/gorm"
)

type VariantFeature struct {
	gorm.Model
	FeatureID uint    `gorm:"primaryKey"`
	VariantID uint    `gorm:"primaryKey"`
	Feature   Feature `gorm:"foreignKey:FeatureID"`
	Value     string
}

func (variantFeature *VariantFeature) SetValue(value any) error {
	valueType, stringValue, err := GetEntityValueAndType(value)
	if err != nil {
		return fmt.Errorf("could not set variant feature value: %w", err)
	}

	// Access the linked variant
	if variantFeature.Feature.ID == 0 || variantFeature.Feature.Type != valueType {
		return fmt.Errorf("assigns invalid value type to feature: %w", err)
	}
	variantFeature.Value = stringValue
	return nil
}

func (variantFeature *VariantFeature) GetValue() (any, error) {
	if variantFeature.Feature.ID == 0 {
		return nil, fmt.Errorf("feature not found")
	}
	value, err := GetPayloadValue(variantFeature.Value, variantFeature.Feature.Type)
	if err != nil {
		return nil, fmt.Errorf("could not get variant feature value: %w", err)
	}
	return value, nil
}

package model

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/util"
	"gorm.io/gorm"
)

type VariantFeature struct {
	gorm.Model
	FeatureID uint    `gorm:"primaryKey"`
	VariantID uint    `gorm:"primaryKey"`
	Feature   Feature `gorm:"foreignKey:FeatureID"`
	Variant   Variant `gorm:"foreignKey:VariantID"`
	Value     string
}

func (variantFeature *VariantFeature) SetValue(value any) *fiber.Error {
	valueType, stringValue, err := util.GetValueTypeAndString(value)
	if err != nil {
		return err
	}

	if variantFeature.Feature.ID == 0 || variantFeature.Feature.Type != valueType {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid feature or value type")
	}
	variantFeature.Value = stringValue
	return nil
}

func (variantFeature *VariantFeature) GetValue() (any, *fiber.Error) {
	if variantFeature.Feature.ID == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Feature not found")
	}
	value, err := util.GetJsonValue(variantFeature.Value, variantFeature.Feature.Type)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "could not get variant feature value: "+err.Error())
	}
	return value, nil
}

type VariantFeatureWritePayload struct {
	FeatureId uint `json:"feature_id" validate:"required"`
	Value     any  `json:"value" validate:"required"`
}

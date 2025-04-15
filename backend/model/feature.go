package model

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/util"
	"gorm.io/gorm"
)

type Feature struct {
	gorm.Model
	Name         string `gorm:"not null;unique"`
	Type         string `gorm:"not null"`
	DefaultValue string `gorm:"not null"`
}

func (feature *Feature) UpdateFromPayload(payload *FeatureWritePayload) error {
	feature.Name = payload.Name
	valueType, defaultValue, err := util.GetValueTypeAndString(payload.Value)
	if err != nil {
		return err
	}
	feature.Type = valueType
	feature.DefaultValue = defaultValue
	return nil
}

func NewFeature(feature *FeatureWritePayload) (*Feature, error) {
	valueType, defaultValue, err := util.GetValueTypeAndString(feature.Value)
	if err != nil {
		return &Feature{}, errors.New("Could not instantiate feature entity: " + err.Error())
	}
	return &Feature{
		Name:         feature.Name,
		Type:         valueType,
		DefaultValue: defaultValue,
	}, nil
}

type FeaturePayload struct {
	Id uint `json:"id"`
	FeatureWritePayload
}

type FeatureWritePayload struct {
	Name  string `json:"name" validate:"required"`
	Value any    `json:"value" validate:"required"`
}

func NewFeaturePayload(feature *Feature) (*FeaturePayload, *fiber.Error) {
	value, err := util.GetJsonValue(feature.DefaultValue, feature.Type)

	if err != nil {
		return &FeaturePayload{}, fiber.NewError(fiber.StatusInternalServerError, "Could not instantiate feature payload: "+err.Error())
	}

	return &FeaturePayload{
		Id: feature.ID,
		FeatureWritePayload: FeatureWritePayload{
			Name:  feature.Name,
			Value: value,
		},
	}, nil
}

type FeatureInfo struct {
	VariantId   uint   `json:"variant_id"`
	VariantName string `json:"variant_name"`
	VariantSize uint   `json:"variant_size"`
	TestId      uint   `json:"test_id"`
	TestName    string `json:"test_name"`
	Value       any    `json:"value"`
}

type FeatureSet map[string]FeatureInfo

package model

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Variant struct {
	gorm.Model
	Name     string `gorm:"not null;primary_key"`
	Size     int    `gorm:"not null"`
	Test     Test
	TestID   uint      `gorm:"not null;primary_key"`
	Features []Feature `gorm:"many2many:variant_features;foreignKey:ID;joinForeignKey:VariantID;References:ID;joinReferences:FeatureID"`
}

func (variant *Variant) UpdateFromPayload(payload VariantWritePayload) error {
	variant.Name = payload.Name
	variant.Size = payload.Size
	return nil
}

func NewVariant(payload VariantWritePayload) Variant {
	return Variant{
		Name:     payload.Name,
		Size:     payload.Size,
		Features: make([]Feature, 0),
	}
}

type VariantWritePayload struct {
	Name string `json:"name" validate:"required"`
	Size int    `json:"size" validate:"required,min=5,max=50"`
}

type VariantPayload struct {
	VariantWritePayload
	Id       uint             `json:"id"`
	Features []FeaturePayload `json:"features"`
}

func NewVariantPayload(variant Variant) (VariantPayload, error) {
	features := make([]FeaturePayload, len(variant.Features))
	for i, feature := range variant.Features {
		payload, err := NewFeaturePayload(&feature)
		if err != nil {
			return VariantPayload{}, fiber.NewError(fiber.StatusInternalServerError, "Could not convert feature entity to payload: "+err.Error())
		}
		features[i] = *payload
	}
	return VariantPayload{
		VariantWritePayload: VariantWritePayload{
			Name: variant.Name,
		},
		Id:       variant.ID,
		Features: features,
	}, nil
}

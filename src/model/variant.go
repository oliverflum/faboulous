package model

import (
	"gorm.io/gorm"
)

type Variant struct {
	gorm.Model
	Name     string `gorm:"not null;primary_key"`
	Test     Test
	TestID   uint      `gorm:"not null;primary_key"`
	Features []Feature `gorm:"many2many:variant_features;foreignKey:ID;joinForeignKey:VariantID;References:ID;joinReferences:FeatureID"`
}

func (variant *Variant) UpdateFromPayload(payload VariantPayload) error {
	variant.Name = payload.Name
	return nil
}

func NewVariantEntity(payload VariantPayload) Variant {
	return Variant{
		Name:     payload.Name,
		Features: make([]Feature, 0),
	}
}

func NewVariantPayload(entity Variant) VariantPayload {
	features := make([]FeaturePayload, len(entity.Features))
	for i, feature := range entity.Features {
		payload, err := NewFeaturePayload(feature)
		if err != nil {
			panic("Could not convert feature entity to payload: " + err.Error())
		}
		features[i] = payload
	}
	return VariantPayload{
		Id:       entity.ID,
		Name:     entity.Name,
		Features: features,
	}
}

type VariantPayload struct {
	Id       uint             `json:"id"`
	Name     string           `json:"name" validate:"required"`
	Features []FeaturePayload `json:"features"`
}

package model

import (
	"gorm.io/gorm"
)

type VariantEntity struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Test     TestEntity
	TestID   uint
	Features []FeatureEntity `gorm:"many2many:variant_features;"`
}

func (variant *VariantEntity) TableName() string {
	return "variant"
}

func NewVariantEntity(payload VariantPayload) VariantEntity {
	features := make([]FeatureEntity, len(payload.Features))
	for i, feature := range payload.Features {
		entity, err := NewFeatureEntity(feature)
		if err != nil {
			panic("Could not convert feature payload to entity: " + err.Error())
		}
		features[i] = entity
	}
	return VariantEntity{
		Name:     payload.Name,
		Features: features,
	}
}

func NewVariantPayload(entity VariantEntity) VariantPayload {
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
	Features []FeaturePayload `json:"features" validate:"required"`
}

package model

import (
	"gorm.io/gorm"
)

type VariantEntity struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Test     TestEntity
	TestID   uint
	Features []FeaturePayload `gorm:"many2many:variant_feature;"`
}

func (variant *VariantEntity) TableName() string {
	return "variant"
}

func NewVariantEntity(payload VariantPayload) VariantEntity {
	return VariantEntity{
		Name:     payload.Name,
		Features: payload.Features,
	}
}

func NewVariantPayload(entity VariantEntity) VariantPayload {
	return VariantPayload{
		Id:       entity.ID,
		Name:     entity.Name,
		Features: entity.Features,
	}
}

type VariantPayload struct {
	Id       uint             `json:"id" validate:"required"`
	Name     string           `json:"name" validate:"required"`
	Features []FeaturePayload `json:"features" validate:"required"`
}

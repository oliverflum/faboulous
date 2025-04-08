package model

import (
	"gorm.io/gorm"
)

type VariantEntity struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Test     Test
	TestID   uint
	Features []FeaturePayload `gorm:"many2many:variant_feature;"`
}

func (variant *VariantEntity) TableName() string {
	return "variant"
}

type VariantPayload struct {
	Id       uint             `json:"id" validate:"required"`
	Name     string           `json:"name" validate:"required"`
	Features []FeaturePayload `json:"features" validate:"required"`
}

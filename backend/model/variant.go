package model

import (
	"gorm.io/gorm"
)

type Variant struct {
	gorm.Model
	Name     string `gorm:"primary_key"`
	TestID   uint   `gorm:"primary_key"`
	Size     uint   `gorm:"not null"`
	Test     Test
	Features []Feature `gorm:"many2many:variant_features;foreignKey:ID;joinForeignKey:VariantID;References:ID;joinReferences:FeatureID"`
}

func (variant *Variant) UpdateFromPayload(payload VariantWritePayload) {
	variant.Name = payload.Name
	variant.Size = payload.Size
}

type VariantWritePayload struct {
	Name string `json:"name" validate:"required"`
	Size uint   `json:"size" validate:"required,min=5,max=50"`
}

type VariantPayload struct {
	VariantWritePayload
	Id       uint             `json:"id"`
	Features []FeaturePayload `json:"features,omitempty"`
}

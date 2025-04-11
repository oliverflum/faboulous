package model

import (
	"gorm.io/gorm"
)

const (
	HASH   = "HASH"
	RANDOM = "RANDOM"
)

type TestEntity struct {
	gorm.Model
	Name     string          `gorm:"unique;not null"`
	Active   bool            `gorm:"default:false"`
	Method   string          `gorm:"not null"`
	Variants []VariantEntity `gorm:"foreignKey:TestID"`
}

func NewTestEntity(payload TestPayload) TestEntity {
	return TestEntity{
		Name:     payload.Name,
		Active:   payload.Active,
		Method:   payload.Method,
		Variants: make([]VariantEntity, 0),
	}
}

type TestPayload struct {
	Id       uint             `json:"id"`
	Name     string           `json:"name" validate:"required"`
	Active   bool             `json:"active"`
	Method   string           `json:"method" validate:"required,oneof=HASH RANDOM"`
	Variants []VariantPayload `json:"variants"`
}

func NewTestPayload(entity TestEntity) TestPayload {
	variants := make([]VariantPayload, len(entity.Variants))
	for i, variantEntity := range entity.Variants {
		variants[i] = NewVariantPayload(variantEntity)
	}
	return TestPayload{
		Id:       entity.ID,
		Name:     entity.Name,
		Active:   entity.Active,
		Method:   entity.Method,
		Variants: variants,
	}
}

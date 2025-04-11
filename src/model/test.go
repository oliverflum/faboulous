package model

import (
	"gorm.io/gorm"
)

const (
	HASH   = "HASH"
	RANDOM = "RANDOM"
)

type Test struct {
	gorm.Model
	Name     string    `gorm:"unique;not null"`
	Active   bool      `gorm:"default:false"`
	Method   string    `gorm:"not null"`
	Variants []Variant `gorm:"foreignKey:TestID"`
}

func NewTestEntity(payload TestPayload) Test {
	return Test{
		Name:     payload.Name,
		Active:   payload.Active,
		Method:   payload.Method,
		Variants: make([]Variant, 0),
	}
}

type TestPayload struct {
	Id       uint             `json:"id"`
	Name     string           `json:"name" validate:"required"`
	Active   bool             `json:"active"`
	Method   string           `json:"method" validate:"required,oneof=HASH RANDOM"`
	Variants []VariantPayload `json:"variants"`
}

func NewTestPayload(test Test) TestPayload {
	variants := make([]VariantPayload, len(test.Variants))
	for i, variant := range test.Variants {
		variants[i] = NewVariantPayload(variant)
	}
	return TestPayload{
		Id:       test.ID,
		Name:     test.Name,
		Active:   test.Active,
		Method:   test.Method,
		Variants: variants,
	}
}

package model

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type Method string

const (
	HASH   Method = "HASH"
	RANDOM Method = "RANDOM"
)

func (self *Method) Scan(method interface{}) error {
	*self = Method(method.([]byte))
	return nil
}

func (self Method) Value() (driver.Value, error) {
	return string(self), nil
}

type TestEntity struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Active   bool
	Method   Method
	Variants []VariantEntity
}

func NewTestEntity(payload TestPayload) TestEntity {
	variants := make([]VariantEntity, len(payload.Variants))
	for i, variantPayload := range payload.Variants {
		variants[i] = NewVariantEntity(variantPayload)
	}
	return TestEntity{
		Name:     payload.Name,
		Active:   payload.Active,
		Method:   payload.Method,
		Variants: variants,
	}
}

type TestPayload struct {
	Id       uint             `json:"id" validate:"required"`
	Name     string           `json:"name" validate:"required"`
	Active   bool             `json:"active" validate:"required"`
	Method   Method           `json:"method" validate:"required, oneof HASH RANDOM"`
	Variants []VariantPayload `json:"variants" validate:"required"`
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

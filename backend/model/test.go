package model

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	HASH   = "HASH"
	RANDOM = "RANDOM"
)

type Test struct {
	gorm.Model
	Name                  string    `gorm:"unique;not null"`
	Active                bool      `gorm:"default:false"`
	Method                string    `gorm:"not null"`
	CollapseControlGroups bool      `gorm:"default:true"`
	Variants              []Variant `gorm:"foreignKey:TestID"`
}

func (t *Test) UpdateFromPayload(payload *TestWritePayload) error {
	t.Name = payload.Name
	t.Active = payload.Active
	t.Method = payload.Method
	return nil
}

func (t *Test) AppendVariants(db *gorm.DB, payload *TestPayload) error {
	// Iterate over variants and check if they exist ind db
	for _, variant := range payload.Variants {
		var existingVariant Variant
		result := db.First(&existingVariant, variant.Id)
		if result.RowsAffected == 0 {
			return errors.New("Variant with id " + fmt.Sprint(variant.Id) + " does not exist")
		}
		t.Variants = append(t.Variants, existingVariant)
	}
	return nil
}

func NewTest(payload *TestWritePayload) Test {
	return Test{
		Name:                  payload.Name,
		Active:                payload.Active,
		Method:                payload.Method,
		CollapseControlGroups: payload.CollapseControlGroups,
		Variants:              make([]Variant, 0),
	}
}

type TestWritePayload struct {
	Name                  string `json:"name" validate:"required"`
	Active                bool   `json:"active"`
	Method                string `json:"method" validate:"required,oneof=HASH RANDOM"`
	CollapseControlGroups bool   `json:"collapseControlGroups"`
}
type TestPayload struct {
	TestWritePayload
	Id       uint             `json:"id"`
	Variants []VariantPayload `json:"variants"`
}

func NewTestPayload(test *Test) (TestPayload, error) {
	variants := make([]VariantPayload, len(test.Variants))
	for i, variant := range test.Variants {
		payload, err := NewVariantPayload(variant)
		if err != nil {
			return TestPayload{}, fiber.NewError(fiber.StatusInternalServerError, "Could not convert variant entity to payload: "+err.Error())
		}
		variants[i] = payload
	}
	return TestPayload{
		TestWritePayload: TestWritePayload{
			Name:                  test.Name,
			Active:                test.Active,
			Method:                test.Method,
			CollapseControlGroups: test.CollapseControlGroups,
		},
		Id:       test.ID,
		Variants: variants,
	}, nil
}

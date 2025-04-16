package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const (
	HASH   = "HASH"
	RANDOM = "RANDOM"
)

type Test struct {
	gorm.Model
	Name                    string    `gorm:"unique;not null"`
	Active                  bool      `gorm:"default:false"`
	Method                  string    `gorm:"not null"`
	CollapseControlVariants bool      `gorm:"default:true"`
	Variants                []Variant `gorm:"foreignKey:TestID;constraint:OnDelete:CASCADE"`
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

type TestWritePayload struct {
	Name                    string `json:"name" validate:"required"`
	Active                  bool   `json:"active"`
	Method                  string `json:"method" validate:"required,oneof=HASH RANDOM"`
	CollapseControlVariants bool   `json:"CollapseControlVariants"`
}
type TestPayload struct {
	TestWritePayload
	Id       uint             `json:"id"`
	Variants []VariantPayload `json:"variants,omitempty"`
}

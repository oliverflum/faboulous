package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/util"
	"gorm.io/gorm"
)

// sendVariantResponse handles the common logic for sending variant responses
func SendVariantResponse(c *fiber.Ctx, variant model.Variant, statusCode int) error {
	payload, err := model.NewVariantPayload(variant, db.GetDB())
	if err != nil {
		return err
	}
	return c.Status(statusCode).JSON(payload)
}

// getVariantByID retrieves a variant by ID and test ID, returns an error if not found
func GetVariant(variantID uint, preloadFeatures bool) (*model.Variant, *fiber.Error) {
	var variant model.Variant
	var result *gorm.DB
	if preloadFeatures {
		result = db.GetDB().
			Preload("Features").
			Where("id = ?", variantID).
			First(&variant)
	} else {
		result = db.GetDB().Where("id = ?", variantID).First(&variant)
	}

	if result.Error != nil {
		return nil, util.HandleGormError(result)
	}

	return &variant, nil
}

func CheckVariantExists(name string, testID uint) bool {
	var existingVariant model.Variant
	result := db.GetDB().Where("name = ? AND test_id = ?", name, testID).First(&existingVariant)
	return result.RowsAffected > 0
}

func CheckVariantSizeConstraints(db *gorm.DB, test *model.Test, variant *model.Variant, variantSize uint) error {
	variants := test.Variants

	usedUpSize := uint(0)
	biggestVariantSize := uint(0)

	variantExists := false

	for _, v := range variants {
		// Skip the variant we're updating in calculations
		if variant != nil && v.ID == variant.ID {
			variantExists = true
			continue
		}

		usedUpSize += v.Size
		if v.Size > biggestVariantSize {
			biggestVariantSize = v.Size
		}
	}

	variantCount := len(variants)
	if !variantExists {
		variantCount++
	}

	newSize := usedUpSize + variantSize
	if variantSize > biggestVariantSize {
		biggestVariantSize = variantSize
	}

	// Check constraints based on test settings
	if test.CollapseControlVariants {
		if newSize > 100-biggestVariantSize {
			return fiber.NewError(fiber.StatusBadRequest,
				fmt.Sprintf("Total size of all variants (%d) exceeds available space after control variant allocation (%d)",
					newSize, 100-biggestVariantSize))
		}
	} else {
		if newSize*2 > 100 {
			return fiber.NewError(fiber.StatusBadRequest,
				fmt.Sprintf("Total size of all variants and their controls (%d) exceeds 100%%", newSize*2))
		}
	}

	return nil
}

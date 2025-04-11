package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/db"
	"github.com/oliverflum/faboulous/model"
	"github.com/oliverflum/faboulous/util"
	"gorm.io/gorm"
)

func updateVariantFeatureValue(tx *gorm.DB, variant *model.Variant, feature *model.FeaturePayload, value any) error {

	var existingFeature model.Feature
	res := tx.First(&existingFeature, "name = ?", feature.Name)
	if res.Error != nil {
		return util.HandleGormError(res.Error)
	}

	valueType, stringValue, err := model.GetEntityValueAndType(value)
	if err != nil {
		return err
	}

	if existingFeature.Type != valueType {
		return util.FabolousError{
			Code:    fiber.StatusBadRequest,
			Message: "value type mismatch: " + existingFeature.Type + " != " + valueType,
		}
	}

	var variantFeature *model.VariantFeature
	res = tx.Find(&variantFeature, "variant_id = ? AND feature_id = ?", variant.ID, existingFeature.ID)
	if res.RowsAffected == 0 {
		variantFeature = new(model.VariantFeature)
		variantFeature.FeatureID = existingFeature.ID
		variantFeature.VariantID = variant.ID
	}

	variantFeature.Value = stringValue

	result := tx.Save(&variantFeature)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func AddVariant(c *fiber.Ctx) error {
	testID := c.Params("testId")
	if testID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Test ID is required")
	}

	var payload model.VariantPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	errors := util.ValidateStruct(payload)
	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Check if test exists
	var test model.Test
	result := db.GetDB().First(&test, testID)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Test not found")
	}

	// Check if variant with same name already exists for this test
	var existingVariant model.Variant
	result = db.GetDB().Where("name = ? AND test_id = ?", payload.Name, testID).First(&existingVariant)
	if result.RowsAffected > 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Variant with this name already exists for this test")
	}

	variant := model.NewVariantEntity(payload)

	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		variant.TestID = test.ID

		result = tx.Create(&variant)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}

		for _, feature := range payload.Features {
			err := updateVariantFeatureValue(tx, &variant, &feature, feature.Value)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return util.SendErrorRes(c, err)
	}

	// Load the variant with its features
	result = db.GetDB().Preload("Features").First(&variant, variant.ID)
	if result.Error != nil {
		return util.SendErrorRes(c, result.Error)
	}

	return c.Status(fiber.StatusOK).JSON(model.NewVariantPayload(variant))
}

func UpdateVariant(c *fiber.Ctx) error {
	testID := c.Params("testId")
	variantID := c.Params("id")
	if testID == "" || variantID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Test ID and Variant ID are required")
	}

	var payload model.VariantPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Validate the payload
	errors := util.ValidateStruct(payload)
	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Check if variant exists and belongs to the test
	var variant model.Variant
	result := db.GetDB().Where("id = ? AND test_id = ?", variantID, testID).First(&variant)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Variant not found")
	}

	// Update variant
	err := variant.UpdateFromPayload(payload)
	if err != nil {
		return util.SendErrorRes(c, err)
	}

	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		result = tx.Save(&variant)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}

		for _, feature := range payload.Features {
			err := updateVariantFeatureValue(tx, &variant, &feature, feature.Value)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return util.SendErrorRes(c, err)
	}

	// Load the variant with its features
	result = db.GetDB().Preload("Features").First(&variant, variant.ID)
	if result.Error != nil {
		return util.SendErrorRes(c, result.Error)
	}

	return c.Status(fiber.StatusOK).JSON(model.NewVariantPayload(variant))
}

func DeleteVariant(c *fiber.Ctx) error {
	testID := c.Params("testId")
	variantID := c.Params("id")
	if testID == "" || variantID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Test ID and Variant ID are required")
	}

	// Check if variant exists and belongs to the test
	var variant model.Variant
	result := db.GetDB().Where("id = ? AND test_id = ?", variantID, testID).First(&variant)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Variant not found")
	}

	result = db.GetDB().Delete(&variant)
	if result.Error != nil {
		return util.SendErrorRes(c, util.HandleGormError(result.Error))
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

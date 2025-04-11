package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/model"
	"github.com/oliverflum/faboulous/util"
)

func ListVariants(c *fiber.Ctx) error {
	testID := c.Params("testId")
	if testID == "" {
		return c.Status(400).SendString("Test ID is required")
	}

	var variants []model.VariantEntity
	result := util.GetDB().Where("test_id = ?", testID).Find(&variants)
	if result.Error != nil {
		return c.Status(500).SendString(result.Error.Error())
	}

	// Convert entities to payloads
	variantPayloads := make([]model.VariantPayload, len(variants))
	for i, variant := range variants {
		variantPayloads[i] = model.NewVariantPayload(variant)
	}

	return c.Status(200).JSON(variantPayloads)
}

func AddVariant(c *fiber.Ctx) error {
	testID := c.Params("testId")
	if testID == "" {
		return c.Status(400).SendString("Test ID is required")
	}

	var payload model.VariantPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Validate the payload
	errors := util.ValidateStruct(payload)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	// Check if test exists
	var test model.TestEntity
	result := util.GetDB().First(&test, testID)
	if result.RowsAffected == 0 {
		return c.Status(404).SendString("Test not found")
	}

	// Check if variant with same name already exists for this test
	var existingVariant model.VariantEntity
	result = util.GetDB().Where("name = ? AND test_id = ?", payload.Name, testID).First(&existingVariant)
	if result.RowsAffected > 0 {
		return c.Status(400).SendString("Variant with this name already exists for this test")
	}

	variant := model.NewVariantEntity(payload)
	variant.TestID = test.ID

	result = util.GetDB().Create(&variant)
	if result.Error != nil {
		return c.Status(500).SendString(result.Error.Error())
	}

	return c.Status(201).JSON(model.NewVariantPayload(variant))
}

func UpdateVariant(c *fiber.Ctx) error {
	testID := c.Params("testId")
	variantID := c.Params("id")
	if testID == "" || variantID == "" {
		return c.Status(400).SendString("Test ID and Variant ID are required")
	}

	var payload model.VariantPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Validate the payload
	errors := util.ValidateStruct(payload)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	// Check if variant exists and belongs to the test
	var variant model.VariantEntity
	result := util.GetDB().Where("id = ? AND test_id = ?", variantID, testID).First(&variant)
	if result.RowsAffected == 0 {
		return c.Status(404).SendString("Variant not found")
	}

	// Check if new name conflicts with existing variant
	if payload.Name != variant.Name {
		var existingVariant model.VariantEntity
		result = util.GetDB().Where("name = ? AND test_id = ? AND id != ?", payload.Name, testID, variantID).First(&existingVariant)
		if result.RowsAffected > 0 {
			return c.Status(400).SendString("Variant with this name already exists for this test")
		}
	}

	// Update variant
	variant.Name = payload.Name

	// Convert feature payloads to entities
	variant.Features = make([]model.FeatureEntity, 0)
	for _, featurePayload := range payload.Features {
		featureEntity, err := model.NewFeatureEntity(featurePayload)
		if err != nil {
			return c.Status(400).SendString("Invalid feature: " + err.Error())
		}
		variant.Features = append(variant.Features, featureEntity)
	}

	result = util.GetDB().Save(&variant)
	if result.Error != nil {
		return c.Status(500).SendString(result.Error.Error())
	}

	return c.Status(200).JSON(model.NewVariantPayload(variant))
}

func DeleteVariant(c *fiber.Ctx) error {
	testID := c.Params("testId")
	variantID := c.Params("id")
	if testID == "" || variantID == "" {
		return c.Status(400).SendString("Test ID and Variant ID are required")
	}

	// Check if variant exists and belongs to the test
	var variant model.VariantEntity
	result := util.GetDB().Where("id = ? AND test_id = ?", variantID, testID).First(&variant)
	if result.RowsAffected == 0 {
		return c.Status(404).SendString("Variant not found")
	}

	result = util.GetDB().Delete(&variant)
	if result.Error != nil {
		return c.Status(500).SendString(result.Error.Error())
	}

	return c.Status(204).Send(nil)
}

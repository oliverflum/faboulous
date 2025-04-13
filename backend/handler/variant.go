package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/internal/db"
	"github.com/oliverflum/faboulous/backend/internal/util"
	"github.com/oliverflum/faboulous/backend/model"
	"gorm.io/gorm"
)

// sendVariantResponse handles the common logic for sending variant responses
func sendVariantResponse(c *fiber.Ctx, variant model.Variant, statusCode int) error {
	payload, err := model.NewVariantPayload(variant)
	if err != nil {
		return err
	}
	return c.Status(statusCode).JSON(payload)
}

// getVariantByID retrieves a variant by ID and test ID, returns an error if not found
func getVariant(testID uint, variantID uint, preloadFeatures bool) (*model.Variant, error) {
	var variant model.Variant
	var result *gorm.DB
	if preloadFeatures {
		result = db.GetDB().
			Preload("Features").
			Where("id = ? AND test_id = ?", variantID, testID).
			First(&variant)
	} else {
		result = db.GetDB().Where("id = ? AND test_id = ?", variantID, testID).First(&variant)
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Variant not found")
	} else if result.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error fetching variant: "+result.Error.Error())
	}

	return &variant, nil
}

// checkVariantExists checks if a variant with the same name exists for a test
func checkVariantExists(name string, testID uint) error {
	var existingVariant model.Variant
	result := db.GetDB().Where("name = ? AND test_id = ?", name, testID).First(&existingVariant)
	if result.RowsAffected > 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Variant with this name already exists for this test")
	}
	return nil
}

func AddVariant(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}

	payload := &model.VariantWritePayload{}
	err = util.ValidatePayload(c, payload)
	if err != nil {
		return err
	}

	// Check if test exists
	var test model.Test
	result := db.GetDB().First(&test, ids["testId"])
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Test not found")
	}

	if err := checkVariantExists(payload.Name, ids["testId"]); err != nil {
		return err
	}

	variant := model.NewVariant(*payload)
	variant.TestID = test.ID
	result = db.GetDB().Create(&variant)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return sendVariantResponse(c, variant, fiber.StatusCreated)
}

func UpdateVariant(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "id"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}

	payload := &model.VariantWritePayload{}
	err = util.ValidatePayload(c, payload)
	if err != nil {
		return err
	}

	variant, err := getVariant(ids["testId"], ids["id"], false)
	if err != nil {
		return err
	}

	// Update variant
	if err := variant.UpdateFromPayload(*payload); err != nil {
		return util.SendErrorRes(c, err)
	}

	result := db.GetDB().Save(variant)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return sendVariantResponse(c, *variant, fiber.StatusOK)
}

func DeleteVariant(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "id"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}

	variant, err := getVariant(ids["testId"], ids["id"], false)
	if err != nil {
		return err
	}

	result := db.GetDB().Delete(variant)
	if result.Error != nil {
		return util.SendErrorRes(c, util.HandleGormError(result.Error))
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

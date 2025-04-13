package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/internal/db"
	"github.com/oliverflum/faboulous/backend/internal/util"
	"github.com/oliverflum/faboulous/backend/model"
)

func checkEntities(testID uint, variantID uint, featureName string) (*model.Variant, *model.Feature, *model.VariantFeature, error) {
	var variant model.Variant
	result := db.GetDB().Where("id = ? AND test_id = ?", variantID, testID).First(&variant)
	if result.Error != nil {
		return nil, nil, nil, fiber.NewError(fiber.StatusNotFound, "Variant not found")
	}

	var feature model.Feature
	result = db.GetDB().First(&feature, "name = ?", featureName)
	if result.Error != nil {
		return &variant, nil, nil, fiber.NewError(fiber.StatusNotFound, "Feature not found")
	}

	var variantFeature model.VariantFeature
	result = db.GetDB().Preload("Feature").Preload("Variant").Where("variant_id = ? AND feature_id = ?", variantID, feature.ID).First(&variantFeature)
	if result.Error != nil {
		return &variant, &feature, nil, nil
	}

	return &variant, &feature, &variantFeature, nil
}

func AddVariantFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID or variant ID")
	}

	var payload model.FeatureWritePayload
	if err := util.ValidatePayload(c, &payload); err != nil {
		return err
	}

	// Check if variant exists and belongs to the test
	variant, feature, existingVariantFeature, err := checkEntities(ids["testId"], ids["variantId"], payload.Name)
	if err != nil {
		return c.Status(err.(*fiber.Error).Code).SendString(err.Error())
	}

	if existingVariantFeature != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Variant feature already exists")
	}

	valueType, value, err := util.GetValueTypeAndString(payload.Value)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if valueType != feature.Type {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid value type")
	}

	variantFeature := model.VariantFeature{
		VariantID: variant.ID,
		FeatureID: feature.ID,
		Value:     value,
	}

	result := db.GetDB().Create(&variantFeature)
	if result.Error != nil {
		return util.SendErrorRes(c, result.Error)
	}

	res := model.FeaturePayload{
		FeatureWritePayload: payload,
		Id:                  variantFeature.ID,
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func UpdateVariantFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId", "id"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID or variant ID")
	}

	var payload model.FeatureWritePayload
	if err := util.ValidatePayload(c, &payload); err != nil {
		return err
	}

	// Check if entities exist
	_, _, variantFeature, err := checkEntities(ids["testId"], ids["variantId"], payload.Name)
	if err != nil {
		return c.Status(err.(*fiber.Error).Code).SendString(err.Error())
	}

	if variantFeature == nil {
		return c.Status(fiber.StatusNotFound).SendString("Variant feature not found")
	}

	// Update the value
	if err := variantFeature.SetValue(payload.Value); err != nil {
		return util.SendErrorRes(c, err)
	}

	result := db.GetDB().Save(&variantFeature)
	if result.Error != nil {
		return util.SendErrorRes(c, result.Error)
	}

	// Create the response payload
	featurePayload, err := model.NewFeaturePayload(&variantFeature.Feature)
	if err != nil {
		return util.SendErrorRes(c, err)
	}
	featurePayload.Value = payload.Value

	return c.Status(fiber.StatusOK).JSON(featurePayload)
}

func DeleteVariantFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId", "id"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID or variant ID")
	}

	// Check if entities exist using gorm find
	var variantFeature model.VariantFeature
	result := db.GetDB().Where("id = ?", ids["id"]).First(&variantFeature)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).SendString("Variant feature not found")
	}

	result = db.GetDB().Delete(&variantFeature)
	if result.Error != nil {
		return util.SendErrorRes(c, result.Error)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

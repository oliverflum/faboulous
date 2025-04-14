package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/service"
	"github.com/oliverflum/faboulous/backend/util"
)

func AddVariantFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId"})
	if err != nil {
		return err
	}

	payload := &model.FeatureWritePayload{}
	valErr := util.ParseAndValidatePayload(c, payload)
	if valErr != nil {
		return valErr
	}

	// Check if variant exists and belongs to the test
	variant, feature, existingVariantFeature, checkErr := service.CheckEntities(ids["testId"], ids["variantId"], payload.Name)
	if checkErr != nil {
		return checkErr
	}

	if existingVariantFeature != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Variant feature already exists")
	}

	valueType, value, valueErr := util.GetValueTypeAndString(payload.Value)
	if valueErr != nil {
		return valueErr
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
		return util.HandleGormError(result)
	}

	res := model.FeaturePayload{
		FeatureWritePayload: *payload,
		Id:                  variantFeature.ID,
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func UpdateVariantFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId", "variantFeatureId"})
	if err != nil {
		return err
	}

	payload := &model.FeatureWritePayload{}
	valErr := util.ParseAndValidatePayload(c, payload)
	if valErr != nil {
		return valErr
	}

	_, _, variantFeature, checkErr := service.CheckEntities(ids["testId"], ids["variantId"], payload.Name)
	if checkErr != nil {
		return checkErr
	}

	if variantFeature == nil {
		return fiber.NewError(fiber.StatusNotFound, "Variant feature not found")
	}

	// Update the value
	if err := variantFeature.SetValue(payload.Value); err != nil {
		return err
	}

	result := db.GetDB().Save(&variantFeature)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	featurePayload, err := model.NewFeaturePayload(&variantFeature.Feature)
	if err != nil {
		return err
	}
	featurePayload.Value = payload.Value

	return c.Status(fiber.StatusOK).JSON(featurePayload)
}

func DeleteVariantFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId", "variantFeatureId"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID or variant ID")
	}

	// Check if entities exist using gorm find
	var variantFeature model.VariantFeature
	result := db.GetDB().Where("id = ?", ids["id"]).First(&variantFeature)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	result = db.GetDB().Delete(&variantFeature)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

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

	payload := &model.VariantFeatureWritePayload{}
	valErr := util.ParseAndValidatePayload(c, payload)
	if valErr != nil {
		return valErr
	}

	variant, feature, existingVariantFeature, checkErr := service.RetrieveEntities(ids["testId"], ids["variantId"], payload.FeatureId)
	if checkErr != nil {
		return checkErr
	}

	if existingVariantFeature != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Variant feature already exists")
	}

	err = service.IsFeatureIsUsedInAnotherTest(feature.ID, ids["testId"])
	if err != nil {
		return err
	}

	valueType, stringValue, err := util.GetValueTypeAndString(payload.Value)
	if err != nil {
		return err
	}

	if feature.Type != valueType {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid value type for feature "+feature.Name)
	}

	variantFeature := model.VariantFeature{
		VariantID: variant.ID,
		FeatureID: feature.ID,
		Value:     stringValue,
	}

	if result := db.GetDB().Create(&variantFeature); result.Error != nil {
		return util.HandleGormError(result)
	}

	updatedVariantFeature := &model.VariantFeature{}
	result := db.GetDB().Preload("Feature").First(&updatedVariantFeature, variantFeature.ID)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	resBody := service.NewVariantFeaturePayload(updatedVariantFeature)

	return c.Status(fiber.StatusCreated).JSON(resBody)
}

func UpdateVariantFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId", "variantFeatureId"})
	if err != nil {
		return err
	}

	payload := &model.VariantFeatureWritePayload{}
	valErr := util.ParseAndValidatePayload(c, payload)
	if valErr != nil {
		return valErr
	}

	_, _, variantFeature, checkErr := service.RetrieveEntities(ids["testId"], ids["variantId"], payload.FeatureId)
	if checkErr != nil {
		return checkErr
	}

	if variantFeature == nil {
		return fiber.NewError(fiber.StatusNotFound, "Variant feature not found")
	}

	if err := variantFeature.SetValue(payload.Value); err != nil {
		return err
	}

	result := db.GetDB().Save(&variantFeature)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	resBody := service.NewVariantFeaturePayload(variantFeature)

	return c.Status(fiber.StatusOK).JSON(resBody)
}

func DeleteVariantFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId", "variantFeatureId"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID or variant ID")
	}

	// Check if entities exist using gorm find
	var variantFeature model.VariantFeature
	result := db.GetDB().Where("id = ?", ids["variantFeatureId"]).First(&variantFeature)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	result = db.GetDB().Unscoped().Delete(&variantFeature)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

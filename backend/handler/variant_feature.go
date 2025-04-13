package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/internal/db"
	"github.com/oliverflum/faboulous/backend/internal/util"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/service"
)

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
	variant, feature, existingVariantFeature, err := service.CheckEntities(ids["testId"], ids["variantId"], payload.Name)
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
	_, _, variantFeature, err := service.CheckEntities(ids["testId"], ids["variantId"], payload.Name)
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

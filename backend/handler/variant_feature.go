package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/service"
	"github.com/oliverflum/faboulous/backend/util"
)

func SetVariantFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId", "featureId"})
	if err != nil {
		return err
	}

	payload := &model.VariantFeatureWritePayload{}
	valErr := util.ParseAndValidatePayload(c, payload)
	if valErr != nil {
		return valErr
	}

	err = service.CheckFeatureUsedInAnotherTest(ids["featureId"], ids["testId"])
	if err != nil {
		return err
	}

	variant, feature, variantFeature, checkErr := service.RetrieveRelatedEntities(ids["testId"], ids["variantId"], ids["featureId"])
	if checkErr != nil {
		return checkErr
	}

	if variantFeature == nil {
		variantFeature = &model.VariantFeature{
			VariantID: variant.ID,
			FeatureID: feature.ID,
			Variant:   *variant,
			Feature:   *feature,
		}
	}

	if err := variantFeature.SetValue(payload.Value); err != nil {
		return err
	}

	result := db.GetDB().Save(&variantFeature)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	resBody, err := service.NewVariantFeaturePayload(variantFeature)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(resBody)
}

func DeleteVariantFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId", "variantFeatureId"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID or variant ID")
	}

	_, _, variantFeature, checkErr := service.RetrieveRelatedEntities(ids["testId"], ids["variantId"], ids["variantFeatureId"])
	if checkErr != nil {
		return checkErr
	} else if variantFeature == nil {
		return fiber.NewError(fiber.StatusNotFound, "Variant feature not found")
	}

	result := db.GetDB().Unscoped().Delete(&variantFeature)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

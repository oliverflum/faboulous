package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/service"
	"github.com/oliverflum/faboulous/backend/util"
)

func ListFeatures(c *fiber.Ctx) error {
	features, err := service.GetAllFeatures()
	if err != nil {
		return err
	}

	if len(features) == 0 {
		return c.Status(fiber.StatusOK).JSON([]model.FeaturePayload{})
	}

	// Convert the features to FeaturePayload
	featurePayloads := make([]*model.FeaturePayload, len(features))
	for i, feature := range features {
		featurePayload, err := service.NewFeaturePayload(feature)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		featurePayloads[i] = featurePayload
	}
	return c.Status(fiber.StatusOK).JSON(featurePayloads)
}

func AddFeature(c *fiber.Ctx) error {
	payload := &model.FeatureWritePayload{}
	valErr := util.ParseAndValidatePayload(c, payload)
	if valErr != nil {
		return c.Status(valErr.Code).SendString(valErr.Message)
	}

	if service.CheckFeatureExists(payload.Name) {
		return c.Status(fiber.StatusBadRequest).SendString("Feature with this name already exists")
	}

	feature, err := service.NewFeature(payload)
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	result := db.GetDB().Create(&feature)
	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return service.SendFeatureResponse(c, feature, fiber.StatusCreated)
}

func GetFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"featureId"})
	if err != nil {
		return err
	}
	feature, err := service.FindFeatureByID(ids["featureId"])
	if err != nil {
		return err
	}

	return service.SendFeatureResponse(c, feature, fiber.StatusOK)
}

func DeleteFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"featureId"})
	if err != nil {
		return err
	}
	feature, err := service.FindFeatureByID(ids["featureId"])
	if err != nil {
		return err
	}

	result := db.GetDB().Unscoped().Delete(feature)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func UpdateFeature(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"featureId"})
	if err != nil {
		return err
	}
	feature, err := service.FindFeatureByID(ids["featureId"])
	if err != nil {
		return err
	}
	payload := &model.FeatureWritePayload{}
	valErr := util.ParseAndValidatePayload(c, payload)
	if valErr != nil {
		return c.Status(valErr.Code).SendString(valErr.Message)
	}

	feature.UpdateFromPayload(payload)
	result := db.GetDB().Save(feature)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return service.SendFeatureResponse(c, feature, fiber.StatusOK)
}

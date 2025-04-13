package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/internal/db"
	"github.com/oliverflum/faboulous/backend/internal/util"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/service"
)

func ListFeatures(c *fiber.Ctx) error {
	var features []*model.Feature

	if err := db.GetDB().Find(&features).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if len(features) == 0 {
		return c.Status(fiber.StatusOK).JSON([]model.FeaturePayload{})
	}

	// Convert the features to FeaturePayload
	featurePayloads := make([]*model.FeaturePayload, len(features))
	for i, feature := range features {
		featurePayload, err := model.NewFeaturePayload(feature)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		featurePayloads[i] = featurePayload
	}
	return c.Status(fiber.StatusOK).JSON(featurePayloads)
}

func AddFeature(c *fiber.Ctx) error {
	payload := &model.FeatureWritePayload{}
	err := util.ValidatePayload(c, payload)
	if err != nil {
		return err
	}

	if service.CheckFeatureExists(payload.Name) {
		return c.Status(fiber.StatusBadRequest).SendString("Feature with this name already exists")
	}

	feature, err := model.NewFeature(payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	result := db.GetDB().Create(&feature)
	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return service.SendFeatureResponse(c, feature, fiber.StatusCreated)
}

func GetFeature(c *fiber.Ctx) error {
	feature, err := service.GetFeatureByID(c.Params("featureId"))
	if err != nil {
		return err
	}

	return service.SendFeatureResponse(c, feature, fiber.StatusOK)
}

func DeleteFeature(c *fiber.Ctx) error {
	feature, err := service.GetFeatureByID(c.Params("featureId"))
	if err != nil {
		return err
	}

	result := db.GetDB().Delete(feature)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func UpdateFeature(c *fiber.Ctx) error {
	payload := &model.FeatureWritePayload{}
	err := util.ValidatePayload(c, payload)
	if err != nil {
		return err
	}

	feature, err := service.GetFeatureByID(c.Params("featureId"))
	if err != nil {
		return err
	}

	feature.UpdateFromPayload(payload)
	result := db.GetDB().Save(feature)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return service.SendFeatureResponse(c, feature, fiber.StatusOK)
}

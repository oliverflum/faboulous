package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/db"
	"github.com/oliverflum/faboulous/model"
	"github.com/oliverflum/faboulous/util"
)

func ListFeatures(c *fiber.Ctx) error {
	var features []model.Feature

	if err := db.GetDB().Find(&features).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if len(features) == 0 {
		return c.Status(fiber.StatusOK).JSON([]model.FeaturePayload{})
	}

	// Convert the features to FeaturePayload
	featurePayloads := make([]model.FeaturePayload, len(features))
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
	featurePayload := new(model.FeaturePayload)

	if err := c.BodyParser(featurePayload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	errors := util.ValidateStruct(*featurePayload)
	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	// Check if feature with same name already exists
	var existingFeature model.Feature
	result := db.GetDB().Where("name = ?", featurePayload.Name).First(&existingFeature)
	if result.RowsAffected > 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Feature with this name already exists")
	}

	feature, err := model.NewFeatureEntity(*featurePayload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	result = db.GetDB().Create(&feature)

	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	resBody, err := model.NewFeaturePayload(feature)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(resBody)
}

func GetFeature(c *fiber.Ctx) error {
	id := c.Params("id")
	var feature model.Feature

	result := db.GetDB().Find(&feature, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	featurePayload, err := model.NewFeaturePayload(feature)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&featurePayload)
}

func DeleteFeature(c *fiber.Ctx) error {
	id := c.Params("id")

	result := db.GetDB().Delete(&model.Feature{}, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func UpdateFeature(c *fiber.Ctx) error {
	id := c.Params("id")

	// Parse the request body into the feature struct
	featurePayload := new(model.FeaturePayload)
	if err := c.BodyParser(&featurePayload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	// Validate the feature payload
	errors := util.ValidateStruct(featurePayload)
	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Check if the feature exists
	var feature model.Feature
	result := db.GetDB().First(&feature, id)
	if result.RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	// Save the updated feature
	feature.UpdateFromPayload(*featurePayload)
	saveResult := db.GetDB().Save(&feature)
	if saveResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(saveResult.Error.Error())
	}

	resBody, err := model.NewFeaturePayload(feature)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&resBody)
}

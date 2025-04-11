package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/model"
	"github.com/oliverflum/faboulous/util"
)

func ListFeatures(c *fiber.Ctx) error {
	var features []model.FeatureEntity

	if err := util.GetDB().Find(&features).Error; err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if len(features) == 0 {
		return c.Status(200).JSON([]model.FeaturePayload{})
	}

	// Convert the features to FeaturePayload
	featurePayloads := make([]model.FeaturePayload, len(features))
	for i, feature := range features {
		featurePayload, err := model.NewFeaturePayload(feature)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		featurePayloads[i] = featurePayload
	}
	return c.Status(200).JSON(featurePayloads)
}

func AddFeature(c *fiber.Ctx) error {
	featurePayload := new(model.FeaturePayload)

	if err := c.BodyParser(featurePayload); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	errors := util.ValidateStruct(*featurePayload)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}
	// Check if feature with same name already exists
	var existingFeature model.FeatureEntity
	result := util.GetDB().Where("name = ?", featurePayload.Name).First(&existingFeature)
	if result.RowsAffected > 0 {
		return c.Status(400).SendString("Feature with this name already exists")
	}

	feature, err := model.NewFeatureEntity(*featurePayload)

	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	result = util.GetDB().Create(&feature)

	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(500).Send(nil)
	}

	resBody, err := model.NewFeaturePayload(feature)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Status(200).JSON(resBody)
}

func GetFeature(c *fiber.Ctx) error {
	id := c.Params("id")
	var feature model.FeatureEntity

	result := util.GetDB().Find(&feature, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	featurePayload, err := model.NewFeaturePayload(feature)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).JSON(&featurePayload)
}

func DeleteFeature(c *fiber.Ctx) error {
	id := c.Params("id")

	result := util.GetDB().Delete(&model.FeatureEntity{}, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(204).Send(nil)
}

func UpdateFeature(c *fiber.Ctx) error {
	id := c.Params("id")

	// Parse the request body into the feature struct
	featurePayload := new(model.FeaturePayload)
	if err := c.BodyParser(&featurePayload); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	// Validate the feature payload
	errors := util.ValidateStruct(featurePayload)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	// Check if the feature exists
	var feature model.FeatureEntity
	result := util.GetDB().First(&feature, id)
	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	// Save the updated feature
	saveResult := util.GetDB().Save(&feature)
	if saveResult.Error != nil {
		return c.Status(500).SendString(saveResult.Error.Error())
	}

	return c.Status(200).JSON(&feature)
}

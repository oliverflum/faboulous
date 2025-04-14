package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/util"
)

// sendFeatureResponse handles the common logic for sending feature responses
func SendFeatureResponse(c *fiber.Ctx, feature *model.Feature, statusCode int) error {
	featurePayload, err := model.NewFeaturePayload(feature)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(statusCode).JSON(featurePayload)
}

// checkFeatureExists checks if a feature exists by name and returns an error if it does
func CheckFeatureExists(name string) bool {
	var existingFeature model.Feature
	result := db.GetDB().Where("name = ?", name).First(&existingFeature)
	return result.RowsAffected > 0
}

// getFeatureByID retrieves a feature by ID and returns an error if not found
func GetFeatureByID(id uint) (*model.Feature, *fiber.Error) {
	var feature model.Feature
	result := db.GetDB().First(&feature, "id = ?", id)
	if result.Error != nil {
		return nil, util.HandleGormError(result)
	}
	return &feature, nil
}

func GetAllFeatures() ([]*model.Feature, *fiber.Error) {
	var features []*model.Feature
	result := db.GetDB().Find(&features)
	if result.Error != nil {
		return nil, util.HandleGormError(result)
	}
	return features, nil
}

package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/util"
)

// sendFeatureResponse handles the common logic for sending feature responses
func SendFeatureResponse(c *fiber.Ctx, feature *model.Feature, statusCode int) error {
	featurePayload, err := NewFeaturePayload(feature)
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
func FindFeatureByID(id uint) (*model.Feature, *fiber.Error) {
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

func NewFeaturePayload(feature *model.Feature) (*model.FeaturePayload, *fiber.Error) {
	value, err := util.GetJsonValue(feature.DefaultValue, feature.Type)

	if err != nil {
		return &model.FeaturePayload{}, fiber.NewError(fiber.StatusInternalServerError, "Could not instantiate feature payload: "+err.Error())
	}

	return &model.FeaturePayload{
		Id: feature.ID,
		FeatureWritePayload: model.FeatureWritePayload{
			Name:  feature.Name,
			Value: value,
		},
	}, nil
}

func NewFeature(feature *model.FeatureWritePayload) (*model.Feature, *fiber.Error) {
	valueType, defaultValue, err := util.GetValueTypeAndString(feature.Value)
	if err != nil {
		return &model.Feature{}, fiber.NewError(fiber.StatusInternalServerError, "Could not instantiate feature entity: "+err.Error())
	}
	return &model.Feature{
		Name:         feature.Name,
		Type:         valueType,
		DefaultValue: defaultValue,
	}, nil
}

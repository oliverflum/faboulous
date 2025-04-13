package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/internal/db"
	"github.com/oliverflum/faboulous/backend/model"
)

func CheckEntities(testID uint, variantID uint, featureName string) (*model.Variant, *model.Feature, *model.VariantFeature, error) {
	var variant model.Variant
	result := db.GetDB().Where("id = ? AND test_id = ?", variantID, testID).First(&variant)
	if result.Error != nil {
		return nil, nil, nil, fiber.NewError(fiber.StatusNotFound, "Variant not found")
	}

	var feature model.Feature
	result = db.GetDB().First(&feature, "name = ?", featureName)
	if result.Error != nil {
		return &variant, nil, nil, fiber.NewError(fiber.StatusNotFound, "Feature not found")
	}

	var variantFeature model.VariantFeature
	result = db.GetDB().Preload("Feature").Preload("Variant").Where("variant_id = ? AND feature_id = ?", variantID, feature.ID).First(&variantFeature)
	if result.Error != nil {
		return &variant, &feature, nil, nil
	}

	return &variant, &feature, &variantFeature, nil
}

func GetVariantFeaturePayload(variantFeature *model.VariantFeature) (*model.FeaturePayload, error) {
	return &model.FeaturePayload{
		Id: variantFeature.ID,
		FeatureWritePayload: model.FeatureWritePayload{
			Name:  variantFeature.Feature.Name,
			Value: variantFeature.Value,
		},
	}, nil
}

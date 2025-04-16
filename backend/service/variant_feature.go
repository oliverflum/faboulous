package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
)

func RetrieveEntities(testID uint, variantID uint, featureId uint) (*model.Variant, *model.Feature, *model.VariantFeature, *fiber.Error) {
	var variant model.Variant
	result := db.GetDB().Where("id = ? AND test_id = ?", variantID, testID).First(&variant)
	if result.Error != nil {
		return nil, nil, nil, fiber.NewError(fiber.StatusNotFound, "Variant not found")
	}

	var feature model.Feature
	result = db.GetDB().First(&feature, featureId)
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

func IsFeatureIsUsedInAnotherTest(featureId uint, testId uint) *fiber.Error {
	inner := db.GetDB().Model(&model.VariantFeature{}).Select("variants.test_id").Joins("left join variants on variant_features.variant_id = variants.id").Where("feature_id = ?", featureId)
	result := db.GetDB().Model(&model.Test{}).Joins("join (?) i on tests.id = i.test_id", inner).Where("tests.id != ?", testId).Scan(&model.Test{})
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check if feature availability")
	} else if result.RowsAffected > 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Feature is already used in another test")
	}
	return nil
}

func NewVariantFeaturePayload(variantFeature *model.VariantFeature) *model.FeaturePayload {
	return &model.FeaturePayload{
		Id: variantFeature.ID,
		FeatureWritePayload: model.FeatureWritePayload{
			Name:  variantFeature.Feature.Name,
			Value: variantFeature.Value,
		},
	}
}

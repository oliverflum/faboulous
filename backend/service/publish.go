package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/util"
)

func createBaseFeatureSet(allFeatures []*model.Feature) *model.FeatureSet {
	featureSet := make(model.FeatureSet)
	for _, feature := range allFeatures {
		featureSet[feature.Name] = feature.DefaultValue
	}
	return &featureSet
}

func updateFeatureSet(featureSet model.FeatureSet, variantId uint) (*model.FeatureSet, *fiber.Error) {
	db := db.GetDB()
	variant := &model.Variant{}
	res := db.Preload("Features").First(variant, variantId)
	if res.Error != nil {
		return nil, util.HandleGormError(res)
	}

	for _, vf := range variant.Features {
		variantFeature := &model.VariantFeature{}
		res := db.First(variantFeature, "variant_id = ? AND feature_id = ?", variant.ID, vf.ID)
		if res.Error != nil {
			return nil, util.HandleGormError(res)
		}
		value, err := util.GetJsonValue(variantFeature.Value, variantFeature.Feature.Type)
		if err != nil {
			return nil, err
		}
		featureSet[vf.Name] = value
	}
	return &featureSet, nil
}

func PublishConfig() error {
	tests, err := GetAllTests(true)
	if err != nil {
		return err
	}

	features, err := GetAllFeatures()
	if err != nil {
		return err
	}

	featureSet := createBaseFeatureSet(features)
	featureTestArrMap := make(db.FeatureTestArrMap)
	for _, test := range tests {
		for _, variant := range test.Variants {
			variantFeatureSet, err := updateFeatureSet(*featureSet, variant.ID)
			if err != nil {
				return err
			}
			featureInfo := model.FeatureInfo{
				VariantId:   variant.ID,
				VariantName: variant.Name,
				VariantSize: variant.Size,
				FeatureSet:  variantFeatureSet,
			}
			featureTestArrMap[test.ID] = append(featureTestArrMap[test.ID], &featureInfo)
		}
	}
	db.SetTestConfigs(featureTestArrMap)
	return nil
}

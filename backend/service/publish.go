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
		value, err := util.GetJsonValue(variantFeature.Value, vf.Type)
		if err != nil {
			return nil, err
		}
		featureSet[vf.Name] = value
	}
	return &featureSet, nil
}

func createTestConfigMap(tests []*model.Test, featureSet *model.FeatureSet) (*db.TestConfigMap, *fiber.Error) {
	newTestConfigMap := make(db.TestConfigMap)
	for _, test := range tests {
		if !test.Active {
			continue
		}
		k := uint(0)
		newTestConfigMap[test.ID] = &db.TestInfo{
			Method:         test.Method,
			FeatureInfoMap: &db.FeatureInfoMap{},
		}
		// Add variant feature sets
		for _, variant := range test.Variants {
			variantFeatureSet, err := updateFeatureSet(*featureSet, variant.ID)
			if err != nil {
				return nil, err
			}
			featureInfo := model.FeatureInfo{
				VariantId:   variant.ID,
				VariantName: variant.Name,
				VariantSize: variant.Size,
				FeatureSet:  variantFeatureSet,
			}
			for j := uint(0); j < variant.Size; j++ {
				(*newTestConfigMap[test.ID].FeatureInfoMap)[k+j] = &featureInfo
			}
			k += variant.Size
		}
		// Add base feature set
		for i := k; i < 100; i++ {
			featureInfo := model.FeatureInfo{
				VariantId:   0,
				VariantName: "base",
				VariantSize: 100 - k,
				FeatureSet:  featureSet,
			}
			(*newTestConfigMap[test.ID].FeatureInfoMap)[i] = &featureInfo
		}
	}
	return &newTestConfigMap, nil
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

	newTestConfigMap, err := createTestConfigMap(tests, featureSet)
	if err != nil {
		return err
	}

	db.SetTestConfigs(newTestConfigMap)
	return nil
}

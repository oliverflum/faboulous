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
		featureSet[feature.Name] = model.FeatureInfo{
			VariantId:   0,
			VariantName: "base",
			VariantSize: 100,
			TestId:      0,
			TestName:    "base",
			Value:       feature.DefaultValue,
		}
	}
	return &featureSet
}

func updateFeatureSet(featureSet model.FeatureSet, variantId uint) (*model.FeatureSet, *fiber.Error) {
	variant := &model.Variant{}
	res := db.GetDB().Preload("Features").Preload("Test").First(variant, variantId)
	if res.Error != nil {
		return nil, util.HandleGormError(res)
	}
	updatedFeatureSet := make(model.FeatureSet)
	for featureName, featureInfo := range featureSet {
		updatedFeatureSet[featureName] = featureInfo
	}
	for _, vf := range variant.Features {
		variantFeature := &model.VariantFeature{}
		res := db.GetDB().First(variantFeature, "variant_id = ? AND feature_id = ?", variant.ID, vf.ID)
		if res.Error != nil {
			return nil, util.HandleGormError(res)
		}
		value, err := util.GetJsonValue(variantFeature.Value, vf.Type)
		if err != nil {
			return nil, err
		}
		updatedFeatureSet[vf.Name] = model.FeatureInfo{
			VariantId:   variant.ID,
			VariantName: variant.Name,
			VariantSize: variant.Size,
			TestId:      variant.Test.ID,
			TestName:    variant.Test.Name,
			Value:       value,
		}
	}
	return &updatedFeatureSet, nil
}

func createTestConfigMap(tests []*model.Test, baseFeatureSet *model.FeatureSet) (*db.TestConfigMap, *fiber.Error) {
	newTestConfigMap := make(db.TestConfigMap)
	for _, test := range tests {
		if !test.Active {
			continue
		}
		k := uint(0)
		newTestConfigMap[test.ID] = &db.TestInfo{
			Method:      test.Method,
			Name:        test.Name,
			FeatureSets: &db.FeatureSetMap{},
		}

		for _, variant := range test.Variants {
			variantFeatureSet, err := updateFeatureSet(*baseFeatureSet, variant.ID)
			if err != nil {
				return nil, err
			}
			for j := uint(0); j < variant.Size; j++ {
				(*newTestConfigMap[test.ID].FeatureSets)[k+j] = variantFeatureSet
			}
			k += variant.Size
		}

		for _, featureInfo := range *baseFeatureSet {
			featureInfo.VariantSize = 100 - k
		}
		for i := k; i < 100; i++ {
			(*newTestConfigMap[test.ID].FeatureSets)[i] = baseFeatureSet
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

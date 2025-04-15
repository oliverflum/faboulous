package service

import (
	"hash/fnv"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
)

type FeatureInfoSetVal struct {
	VariantId   uint   `json:"variant_id"`
	VariantName string `json:"variant_name"`
	VariantSize uint   `json:"variant_size"`
	TestId      uint   `json:"test_id"`
	TestName    string `json:"test_name"`
	Value       any    `json:"value"`
}
type FeatureInfoSet map[string]*FeatureInfoSetVal

func getFeatureInfoKey(differentiator string, method string) (uint, *fiber.Error) {
	switch method {
	case model.RANDOM:
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)
		randomNumber := random.Intn(100)
		return uint(randomNumber), nil
	case model.HASH:
		h := fnv.New32a()
		h.Write([]byte(differentiator))
		return uint(h.Sum32() % 100), nil
	}
	return 0, fiber.NewError(fiber.StatusInternalServerError, "Invalid method")
}

func GetFeatureSet(differentiator string) (FeatureInfoSet, *fiber.Error) {
	testsConfigs := db.GetTestConfigs()
	featureInfoSet := make(FeatureInfoSet)
	for testId, testConfig := range testsConfigs {
		if differentiator == "" && testConfig.Method != model.RANDOM {
			continue
		}
		featureInfoKey, err := getFeatureInfoKey(differentiator, testConfig.Method)
		if err != nil {
			return nil, err
		}

		featureInfo := (*testConfig.FeatureInfoMap)[featureInfoKey]
		for featureName, featureVal := range *featureInfo.FeatureSet {
			featureInfoSet[featureName] = &FeatureInfoSetVal{
				VariantId:   featureInfo.VariantId,
				VariantName: featureInfo.VariantName,
				VariantSize: featureInfo.VariantSize,
				TestId:      testId,
				TestName:    testConfig.Name,
				Value:       featureVal,
			}
		}
		return featureInfoSet, nil
	}
	return featureInfoSet, nil
}

package service

import (
	"hash/fnv"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
)

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

func GetFeatureSet(differentiator string) (*model.FeatureSet, *fiber.Error) {
	testsConfigs := db.GetTestConfigs()
	globalFeatureSet := make(model.FeatureSet)
	for _, testConfig := range testsConfigs {
		if differentiator == "" && testConfig.Method != model.RANDOM {
			continue
		}
		featureInfoKey, err := getFeatureInfoKey(differentiator, testConfig.Method)
		if err != nil {
			return nil, err
		}

		testFeatureSet := (*testConfig.FeatureSets)[featureInfoKey]
		for featureName, featureInfo := range *testFeatureSet {
			globalFeatureSet[featureName] = featureInfo
		}
	}
	return &globalFeatureSet, nil
}

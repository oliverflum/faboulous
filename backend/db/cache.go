package db

import (
	"sync"

	"github.com/oliverflum/faboulous/backend/model"
)

type FeatureSetMap map[uint]*model.FeatureSet

type TestInfo struct {
	Method      string
	Name        string
	FeatureSets *FeatureSetMap
}
type TestConfigMap map[uint]*TestInfo

var testConfigMap TestConfigMap
var testConfigMutex sync.RWMutex

func GetTestConfigs() TestConfigMap {
	testConfigMutex.RLock()
	defer testConfigMutex.RUnlock()

	return testConfigMap
}

func SetTestConfigs(newTestConfigMap *TestConfigMap) {
	testConfigMutex.Lock()
	defer testConfigMutex.Unlock()

	testConfigMap = *newTestConfigMap
}

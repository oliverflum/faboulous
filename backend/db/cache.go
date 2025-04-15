package db

import (
	"sync"

	"github.com/oliverflum/faboulous/backend/model"
)

type FeatureInfoMap map[uint]*model.FeatureInfo
type TestInfo struct {
	Method         string
	FeatureInfoMap *FeatureInfoMap
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

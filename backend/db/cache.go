package db

import (
	"sync"

	"github.com/oliverflum/faboulous/backend/model"
)

type FeatureInfoMap map[uint]*model.FeatureInfo
type TestConfigMap map[uint]*FeatureInfoMap
type FeatureTestArrMap map[uint][]*model.FeatureInfo

var testConfigMap TestConfigMap

var testConfigMutex sync.RWMutex

func GetTestConfigs(key uint) *TestConfigMap {
	testConfigMutex.RLock()
	defer testConfigMutex.RUnlock()

	return &testConfigMap
}

func SetTestConfigs(infos FeatureTestArrMap) {
	newTestConfigMap := make(TestConfigMap)
	for testId, infos := range infos {
		newTestConfigMap[testId] = &FeatureInfoMap{}
		k := uint(0)
		for _, info := range infos {
			for j := uint(0); j < info.VariantSize; j++ {
				(*newTestConfigMap[testId])[k+j] = info
			}
			k += info.VariantSize
		}
	}

	testConfigMutex.Lock()
	defer testConfigMutex.Unlock()

	testConfigMap = newTestConfigMap
}

package utils

import (
	"fmt"
	"os"
)

func ReadEnvVars(envVarNames []string) map[string]string {
	envVarVals := make(map[string]string)
	for _, envVarName := range envVarNames {
		val, varSet := os.LookupEnv(envVarName)
		if !varSet {
			msg := fmt.Sprintf("Environment Variable %s is not set", envVarName)
			panic(msg)
		}
		envVarVals[envVarName] = val
	}
	return envVarVals
}

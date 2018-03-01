package utils

import "os"

func GetEnvOr(envVar, defaultValue string) string {
	val, ok := os.LookupEnv(envVar)
	if ok {
		return val
	}
	return defaultValue
}

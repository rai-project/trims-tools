package micro

import "os"

func getEnvOr(envVar, defaultValue string) string {
	val, ok := os.LookupEnv(envVar)
	if ok {
		return val
	}

	return defaultValue
}

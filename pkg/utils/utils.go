package utils

import (
	"os"
	"regexp"
)

func GetEnvOr(envVar, defaultValue string) string {
	val, ok := os.LookupEnv(envVar)
	if ok {
		return val
	}
	return defaultValue
}

func MatchOneOf(text string, patterns ...string) bool {
	for _, pattern := range patterns {
		if text == pattern {
			return true
		}
		re := regexp.MustCompile(pattern)
		value := re.FindStringSubmatch(text)
		if len(value) > 0 {
			return true
		}
	}
	return false
}

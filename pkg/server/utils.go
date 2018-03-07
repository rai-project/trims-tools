package server

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	EvictionPolicies = []string{
		"never",
		"lru",
		"lcu",
		"fifo",
		"flush",
		"all",
	}
)

func IsValidEvictionPolicy(policy string) (bool, error) {
	for _, p := range EvictionPolicies {
		if p == policy {
			return true, nil
		}
	}
	return false,
		errors.Errorf("the policy %s is not valid and must be one of %s", policy, strings.Join(EvictionPolicies, ", "))
}

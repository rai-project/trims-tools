package workload

import (
	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/gonum/stat/samplemv"
)

// Pareto
// Zipf
// Uniform
// best case : only once choice
// Worst case : always evict
var Generator rand.Rand
var Distributions0 samplemv.Sampler
var Distributions []distuv.Rander

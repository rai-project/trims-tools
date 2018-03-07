package workload

import (
	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/gonum/stat/samplemv"
)

var Generator rand.Rand
var Distributions0 samplemv.Sampler
var Distributions []distuv.Rander

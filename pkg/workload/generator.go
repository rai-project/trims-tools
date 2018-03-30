package workload

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/assets"

	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/gonum/stat/sampleuv"
)

// Pareto
// Zipf (not done)
// Uniform
// Exp
// Weibull
// Poisson

type Dist interface {
	// distuv.Quantiler
	distuv.RandLogProber
	CDF(float64) float64
	Prob(float64) float64
}

type Generator struct {
	done chan (bool)
	name string
	Dist
}

var (
	ValidDistributions = []string{
		"pareto",
		"zipf",
		"uniform",
		"exp", "exponential",
		"weibull",
		"poisson",
	}
	DefaultParetoParameters      = []float64{1, 1.5}
	DefaultZipfParameters        = []float64{}
	DefaultUniformParameters     = []float64{0, 1}
	DefaultExponentialParameters = []float64{0.5}
	DefaultWeibullParameters     = []float64{1.5, 1}
	DefaultPoissonParameters     = []float64{1}
)

func New(distribution string, params []float64) (*Generator, error) {
	if params == nil {
		params = []float64{}
	}
	var rnd *Generator
	switch strings.ToLower(distribution) {
	case "pareto":
		if len(params) != 2 {
			params = DefaultParetoParameters
		}
		rnd = NewPareto(params[0], params[1])
	case "zipf":
		return nil, errors.New("the zipf distribution is not implemented")
	case "uniform":
		if len(params) != 2 {
			params = DefaultUniformParameters
		}
		rnd = NewUniform(params[0], params[1])
	case "exp", "exponential":
		if len(params) != 1 {
			params = DefaultExponentialParameters
		}
		rnd = NewExponential(params[0])
	case "weibull":
		if len(params) != 2 {
			params = DefaultWeibullParameters
		}
		rnd = NewWeibull(params[0], params[1])
	case "poisson":
		if len(params) != 1 {
			params = DefaultPoissonParameters
		}
		rnd = NewPoisson(params[0])
	default:
		return nil, errors.Errorf("the distribution %s is unknown", distribution)
	}

	if rnd == nil {
		return nil, errors.Errorf("the distribution %s is unknown", distribution)
	}

	rnd.name = strings.ToLower(distribution)

	return rnd, nil
}

func NewUniform(min float64, max float64) *Generator {
	return &Generator{
		Dist: distuv.Uniform{Min: min, Max: max},
	}
}

func NewExponential(rate float64) *Generator {
	return &Generator{
		Dist: distuv.Exponential{Rate: rate},
	}
}

func NewWeibull(k float64, lambda float64) *Generator {
	return &Generator{
		Dist: distuv.Weibull{K: k, Lambda: lambda},
	}
}

func NewPoisson(lambda float64) *Generator {
	return &Generator{
		Dist: distuv.Poisson{Lambda: lambda},
	}
}

func NewZipf(...float64) *Generator {
	panic("not implemented")
	return nil
}

func NewPareto(xm float64, alpha float64) *Generator {
	return &Generator{
		Dist: distuv.Pareto{Xm: xm, Alpha: alpha},
	}
}

type ProposalDist struct {
	Dist
}

func (p ProposalDist) ConditionalRand(y float64) float64 {
	return p.Rand()
}

func (p ProposalDist) ConditionalLogProb(x, y float64) float64 {
	return p.LogProb(x)
}

func (g *Generator) Next(at AliasTable, arry []interface{}) interface{} {
	idx := at.Next()
	return arry[idx]
}

func (g *Generator) Next_0(arry []interface{}) interface{} {

	target := g.Dist
	proposal := ProposalDist{distuv.Uniform{Min: 0, Max: float64(len(arry))}}

	imp := sampleuv.Importance{Target: target, Proposal: proposal}

	nSamples := len(arry)
	x := make([]float64, nSamples)
	weights := make([]float64, nSamples)

	imp.SampleWeighted(x, weights)

	fmt.Println(int(x[0]))

	return int(x[0])
}

func (g *Generator) NextMH(arry []interface{}) interface{} {
	if g.name == "uniform" {
		dist := distuv.Uniform{Min: 0, Max: float64(len(arry))}
		idx := int(dist.Rand())
		fmt.Println(idx)
		return arry[idx]
	}
	n := 1
	burnin := 100
	var initial float64
	// target is the distribution from which we would like to sample.
	target := g.Dist

	// proposal is the proposal distribution. Here, we are choosing
	// a tight Gaussian distribution around the current location. In
	// typical problems, if Sigma is too small, it takes a lot of samples
	// to move around the distribution. If Sigma is too large, it can be hard
	// to find acceptable samples.
	proposal := ProposalDist{distuv.Uniform{Min: 0, Max: float64(len(arry))}}

	samples := make([]float64, n+burnin)

	mh := sampleuv.MetropolisHastings{Initial: initial, Target: target, Proposal: proposal, BurnIn: burnin, Rate: 1}
	mh.Sample(samples)

	samples = samples[burnin:]

	idx := int(samples[0])
	fmt.Println(idx)
	return arry[idx]
}

func makeRange(len int) []float64 {
	res := make([]float64, len)
	for ii := range res {
		res[ii] = float64(ii)
	}
	return res
}

func (g *Generator) Generator(arry []interface{}) <-chan interface{} {
	gen := make(chan interface{}, 10)
	at := NewAlias(g.probs(len(arry)), rand.NewSource(0))

	go func() {
		defer close(gen)
		for {
			select {
			case <-g.done:
				return
			default:
				gen <- g.Next(at, arry)
			}
		}
	}()
	return gen
}

func (g *Generator) probs0(len int) []float64 {
	dist := g.Dist
	pmin := 1.0
	switch g.name {
	case "unifrom":
		pmin = 1.0
	case "pareto":
		pmin = 7.0
	case "zipf":
		pmin = 1.0
	case "exp", "exponential":
		pmin = 3.0
	case "weibull":
		pmin = 4.0
	case "poisson":
		pmin = 4.0
	}
	pmax := pmin
	for ii := pmax; ii < 1000; ii += 0.1 {
		e := dist.CDF(ii)
		if 1-e < 0.001 {
			pmax = ii
			break
		}
	}
	res := make([]float64, len)
	for ii := range res {
		res[ii] = dist.Rand() //float64(ii) * pmax / float64(len))
	}
	total := 0.0
	for _, r := range res {
		total = total + r
	}
	for ii := range res {
		res[ii] = res[ii] / total
	}

	fmt.Println(pmax)
	fmt.Println(res)
	return res
}

func (g *Generator) probs(len int) []float64 {
	dist := g.Dist
	res := make([]float64, len)
	for ii := range res {
		res[ii] = dist.Rand() //float64(ii) * pmax / float64(len))
	}
	total := 0.0
	for _, r := range res {
		total = total + r
	}
	for ii := range res {
		res[ii] = res[ii] / total
	}
	return res
}

func (g *Generator) ModelGenerator(models assets.ModelManifests) <-chan assets.ModelManifest {
	gen := make(chan assets.ModelManifest, 10)
	arry := make([]interface{}, len(models))
	for ii, m := range models {
		arry[ii] = m
	}

	at := NewAlias(g.probs(len(arry)), rand.NewSource(0))

	go func() {
		defer close(gen)
		for {
			select {
			case <-g.done:
				return
			default:
				n := g.Next(at, arry)
				gen <- n.(assets.ModelManifest)
			}
		}
	}()
	return gen
}

func (g *Generator) Wait() {
	<-g.done
}

func (g *Generator) Close() {
	close(g.done)
}

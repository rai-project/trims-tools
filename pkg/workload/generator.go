package workload

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/assets"

	"gonum.org/v1/gonum/stat/distuv"
)

// Pareto
// Zipf (not done)
// Uniform
// Exp
// Weibull
// Poisson

type Generator struct {
	done chan (bool)
	distuv.Rander
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
)

func New(distribution string, params []float64) (*Generator, error) {
	switch strings.ToLower(distribution) {
	case "pareto":
		return NewPareto(params[0], params[1]), nil
	case "zipf":
		return nil, errors.New("the zipf distribution is not implemented")
	case "uniform":
		return NewUniform(params[0], params[1]), nil
	case "exp", "exponential":
		return NewExponential(params[0]), nil
	case "weibull":
		return NewWeibull(params[0], params[1]), nil
	case "poisson":
		return NewPoisson(params[0]), nil
	}

	return nil, errors.Errorf("the distribution %s is unknown", distribution)
}

func NewUniform(min float64, max float64) *Generator {
	return &Generator{
		Rander: distuv.Uniform{Min: min, Max: max},
	}
}

func NewExponential(rate float64) *Generator {
	return &Generator{
		Rander: distuv.Exponential{Rate: rate},
	}
}

func NewWeibull(k float64, lambda float64) *Generator {
	return &Generator{
		Rander: distuv.Weibull{K: k, Lambda: lambda},
	}
}

func NewPoisson(lambda float64) *Generator {
	return &Generator{
		Rander: distuv.Poisson{Lambda: lambda},
	}
}

func NewZipf(...float64) *Generator {
	panic("not implemented")
	return nil
}

func NewPareto(xm float64, alpha float64) *Generator {
	return &Generator{
		Rander: distuv.Pareto{Xm: xm, Alpha: alpha},
	}
}

func (g *Generator) Next(arry []interface{}) interface{} {
	r := g.Rand()
	idx := int(r * float64(len(arry)))
	return arry[idx]
}

func (g *Generator) NextModel(models assets.ModelManifests) assets.ModelManifest {
	n := g.Next([]interface{}{models})
	return n.(assets.ModelManifest)
}

func (g *Generator) Generator(arry []interface{}) <-chan interface{} {
	gen := make(chan interface{}, 2)
	go func() {
		defer close(gen)
		for {
			select {
			case <-g.done:
				return
			default:
				gen <- g.Next(arry)
			}
		}
	}()
	return gen
}

func (g *Generator) ModelGenerator(models assets.ModelManifests) <-chan assets.ModelManifest {
	gen := make(chan assets.ModelManifest, 2)
	go func() {
		defer close(gen)
		for {
			select {
			case <-g.done:
				return
			default:
				gen <- g.NextModel(models)
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

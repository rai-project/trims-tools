package client

import (
	"context"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/rai-project/micro18-tools/pkg/workload"

	"github.com/spf13/cast"
)

type Options struct {
	ctx                     context.Context
	original                bool
	profileIO               bool
	iterationCount          int
	debug                   bool
	eagerInitialize         bool
	eagerInitializeAsync    bool
	postprocess             bool
	modelName               string
	uploadProfile           bool
	modelDistribution       string
	modelDistributionParams []float64
	concurrentRunCount      int
	modelIterationCount     int
	profileMemory           bool
	showProgress            bool
	largeModels             bool
	stderr                  io.Writer
	stdout                  io.Writer
}

type Option func(*Options)

var (
	DefaultOptions = Options{
		ctx:                     context.Background(),
		profileIO:               true,
		original:                false,
		iterationCount:          3,
		debug:                   false,
		eagerInitialize:         true,
		eagerInitializeAsync:    false,
		postprocess:             true,
		modelName:               "all",
		uploadProfile:           true,
		modelDistribution:       "none",
		modelDistributionParams: []float64{},
		concurrentRunCount:      runtime.NumCPU(),
		modelIterationCount:     -1,
		profileMemory:           false,
		showProgress:            true,
		largeModels:             false,
		stderr:                  os.Stderr,
		stdout:                  os.Stdout,
	}
)

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.ctx = ctx
	}
}

func ConcurrentRunCount(ii int) Option {
	return func(o *Options) {
		o.concurrentRunCount = ii
	}
}

func UploadProfile(b bool) Option {
	return func(o *Options) {
		o.uploadProfile = b
	}
}

func LargeModels(b bool) Option {
	return func(o *Options) {
		o.largeModels = b
	}
}

func IterationCount(ii int) Option {
	return func(o *Options) {
		o.iterationCount = ii
	}
}

func ModelIterationCount(ii int) Option {
	return func(o *Options) {
		o.modelIterationCount = ii
	}
}

func OriginalMode(b bool) Option {
	return func(o *Options) {
		o.original = b
	}
}

func ProfileIO(b bool) Option {
	return func(o *Options) {
		o.profileIO = b
	}
}

func DebugMode(b bool) Option {
	return func(o *Options) {
		o.debug = b
	}
}

func PostProcess(b bool) Option {
	return func(o *Options) {
		o.postprocess = b
	}
}

func ShowProgress(b bool) Option {
	return func(o *Options) {
		o.showProgress = b
	}
}

func EagerInitialize(b bool) Option {
	return func(o *Options) {
		o.eagerInitialize = b
	}
}

func EagerInitializeAsync(b bool) Option {
	return func(o *Options) {
		o.eagerInitializeAsync = b
	}
}

func ModelName(n string) Option {
	return func(o *Options) {
		o.modelName = n
	}
}

func ProfileMemory(b bool) Option {
	return func(o *Options) {
		o.profileMemory = b
	}
}

func Stdout(w io.Writer) Option {
	return func(o *Options) {
		o.stdout = w
	}
}

func Stderr(w io.Writer) Option {
	return func(o *Options) {
		o.stderr = w
	}
}

func ModelDistribution(dist, params string) Option {
	return func(o *Options) {
		if strings.ToLower(dist) == "none" {
			o.modelDistribution = "none"
			return
		}
		if !workload.IsValidDistribution(dist) {
			panic(
				"the distribution " +
					dist +
					" is not valid. Please specify one of " +
					strings.Join(workload.ValidDistributions, ",") +
					" distributions",
			)
		}
		o.modelDistribution = dist
		if params == "" {
			return
		}
		for _, e := range strings.Split(params, ",") {
			o.modelDistributionParams = append(o.modelDistributionParams, cast.ToFloat64(e))
		}
	}
}

func WithOptions(opts ...Option) *Options {
	os := DefaultOptions
	for _, o := range opts {
		o(&os)
	}
	return &os
}

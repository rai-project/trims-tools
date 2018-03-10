package client

import (
	"context"
	"strings"

	"github.com/spf13/cast"
)

type Options struct {
	ctx                     context.Context
	original                bool
	iterationCount          int
	debug                   bool
	eagerInitialize         bool
	eagerInitializeAsync    bool
	postprocess             bool
	modelName               string
	uploadProfile           bool
	modelDistribution       string
	modelDistributionParams []float64
}

type Option func(*Options)

var (
	DefaultOptions = Options{
		ctx:                  context.Background(),
		original:             false,
		iterationCount:       3,
		debug:                false,
		eagerInitialize:      false,
		eagerInitializeAsync: false,
		postprocess:          false,
		modelName:            "all",
		uploadProfile:        true,
	}
)

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.ctx = ctx
	}
}

func UploadProfile(b bool) Option {
	return func(o *Options) {
		o.uploadProfile = b
	}
}

func IterationCount(ii int) Option {
	return func(o *Options) {
		o.iterationCount = ii
	}
}

func OriginalMode(b bool) Option {
	return func(o *Options) {
		o.original = b
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

func ModelDistribution(dist, params string) Option {
	return func(o *Options) {
		o.modelDistribution = dist
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

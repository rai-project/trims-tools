package server

import "context"

type Options struct {
	ctx                 context.Context
	debug               bool
	evictionPolicy      string
	modelEstimationRate float32
	memoryPercentage    float32
}

type Option func(*Options)

var (
	DefaultOptions = Options{
		ctx:                 context.Background(),
		debug:               false,
		evictionPolicy:      "lru",
		modelEstimationRate: 1.5,
		memoryPercentage:    0.8,
	}
)

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.ctx = ctx
	}
}

func DebugMode(b bool) Option {
	return func(o *Options) {
		o.debug = b
	}
}

func EvictionPolicy(n string) Option {
	return func(o *Options) {
		o.evictionPolicy = n
	}
}

func ModelEstimationRate(n float32) Option {
	return func(o *Options) {
		o.modelEstimationRate = n
	}
}

func MemoryPercentage(n float32) Option {
	return func(o *Options) {
		o.memoryPercentage = n
	}
}

func WithOptions(opts ...Option) *Options {
	os := DefaultOptions
	for _, o := range opts {
		o(&os)
	}
	return &os
}

package trace

import "context"

type Options struct {
	ctx                  context.Context
	iterationCount       int
	eagerInitialize      bool
	eagerInitializeAsync bool
	modelName            string
}

type Option func(*Options)

var (
	DefaultOptions = Options{
		ctx:                  context.Background(),
		iterationCount:       3,
		eagerInitialize:      false,
		eagerInitializeAsync: false,
		modelName:            "All",
	}
)

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.ctx = ctx
	}
}

func IterationCount(ii int) Option {
	return func(o *Options) {
		o.iterationCount = ii
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

func WithOptions(opts ...Option) *Options {
	os := DefaultOptions
	for _, o := range opts {
		o(&os)
	}
	return &os
}

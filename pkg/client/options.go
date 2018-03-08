package client

import "context"

type Options struct {
	ctx                  context.Context
	iterationCount       int
	debug                bool
	eagerInitialize      bool
	eagerInitializeAsync bool
	postprocess          bool
	modelName            string
	uploadProfile        bool
}

type Option func(*Options)

var (
	DefaultOptions = Options{
		ctx:                  context.Background(),
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

func WithOptions(opts ...Option) *Options {
	os := DefaultOptions
	for _, o := range opts {
		o(&os)
	}
	return &os
}

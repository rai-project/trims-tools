package server

import (
	"context"
	"io"
	"os"
)

type Options struct {
	ctx                 context.Context
	debug               bool
	evictionPolicy      string
	modelEstimationRate float32
	memoryPercentage    float32
	uploadProfile       bool
	stderr              io.Writer
	stdout              io.Writer
	persistCPU          bool
}

type Option func(*Options)

var (
	DefaultOptions = Options{
		ctx:                 context.Background(),
		debug:               false,
		evictionPolicy:      "lru",
		modelEstimationRate: 1.0,
		memoryPercentage:    0.8,
		uploadProfile:       true,
		stderr:              os.Stderr,
		stdout:              os.Stdout,
		persistCPU:          true,
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

func PersistCPU(b bool) Option {
	return func(o *Options) {
		o.persistCPU = b
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

func WithOptions(opts ...Option) *Options {
	os := DefaultOptions
	for _, o := range opts {
		o(&os)
	}
	return &os
}

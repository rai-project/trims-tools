package server

import (
	"context"
	"io"
	"os"
)

type Options struct {
	ctx                        context.Context
	id                         string
	debug                      bool
	evictionPolicy             string
	modelEstimationRate        float32
	memoryPercentage           float32
	uploadProfile              bool
	stderr                     io.Writer
	stdout                     io.Writer
	persistCPU                 bool
	writeProfile               bool
	estimateWithInternalMemory bool
}

type Option func(*Options)

var (
	DefaultOptions = Options{
		ctx:                        context.Background(),
		id:                         "",
		debug:                      false,
		evictionPolicy:             "lru",
		modelEstimationRate:        1.0,
		memoryPercentage:           0.8,
		uploadProfile:              true,
		stderr:                     os.Stderr,
		stdout:                     os.Stdout,
		persistCPU:                 true,
		writeProfile:               false,
		estimateWithInternalMemory: true,
	}
)

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.ctx = ctx
	}
}

func ID(s string) Option {
	return func(o *Options) {
		o.id = s
	}
}

func WriteProfile(b bool) Option {
	return func(o *Options) {
		o.writeProfile = b
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

func EstimateWithInternalMemory(b bool) Option {
	return func(o *Options) {
		o.estimateWithInternalMemory = b
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

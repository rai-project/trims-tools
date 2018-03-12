package server

import (
	"io"
	"os"

	"github.com/rai-project/micro18-tools/pkg/builder"
	"github.com/rai-project/micro18-tools/pkg/watcher"
)

type Server struct {
	output  io.Reader
	builder *builder.Builder
	watcher *watcher.Watcher
	options Options
}

func New(opts ...Option) *Server {
	options := WithOptions(opts...)
	return &Server{
		output:  os.Stdout,
		options: *options,
	}
}

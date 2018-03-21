package server

import (
	"io"
	"os"
	"os/exec"

	"github.com/rai-project/micro18-tools/pkg/builder"
	"github.com/rai-project/micro18-tools/pkg/watcher"
)

type Server struct {
	output  io.Reader
	cmd     *exec.Cmd
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

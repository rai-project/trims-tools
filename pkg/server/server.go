package server

import (
	"io"

	"github.com/rai-project/micro18-tools/pkg/builder"
	"github.com/rai-project/micro18-tools/pkg/watcher"
)

type Server struct {
	path    string
	output  io.Reader
	builder *builder.Builder
	watcher *watcher.Watcher
	options Options
}

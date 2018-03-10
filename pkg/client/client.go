package client

import (
	"io"
	"os"

	"github.com/rai-project/micro18-tools/pkg/builder"
	"github.com/rai-project/micro18-tools/pkg/watcher"
)

type Client struct {
	output  io.Reader
	builder *builder.Builder
	watcher *watcher.Watcher
	options Options
}

func New(opts ...Option) *Client {
	options := WithOptions(opts...)
	return &Client{
		output:  os.Stdout,
		options: *options,
	}
}

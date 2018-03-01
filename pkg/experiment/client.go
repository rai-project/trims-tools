package experiment

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"time"

	"github.com/fsnotify/fsnotify"
	tui "github.com/marcusolsson/tui-go"
	builder "github.com/rai-project/micro18-tools/pkg/builder"
	watcher "github.com/rai-project/micro18-tools/pkg/watcher"
)

type Client struct {
	ctx     context.Context
	path    string
	output  io.Reader
	content bytes.Buffer
	builder *builder.Builder
	watcher *watcher.Watcher
	widget  *OutputWidget
}

func NewClient(ctx context.Context, ui tui.UI) *Client {
	ctx = context.WithValue(ctx, "name", "client")

	outputReader, outputWriter := io.Pipe()

	widget := NewOutputWidget("Client", ui)

	client := &Client{
		output:  outputReader,
		content: bytes.Buffer{},
		widget:  widget,
		builder: &builder.Builder{
			Ctx:      context.WithValue(ctx, "name", "client-builder"),
			Stderr:   outputWriter,
			Stdout:   outputWriter,
			BaseDir:  Config.ClientPath,
			BuildCmd: Config.ClientBuildCmd,
		},
	}
	defer func() {
		go client.listen()
	}()
	client.watcher = watcher.New(
		client.Notify,
		Config.ClientPath,
		Config.ServerPath, // changes to the server should trigger the client to rerun
	)
	return client
}

func (c *Client) listen() (err error) {
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			log.Panic(e)
		}
	}()

	nBytes, nChunks := int64(0), int64(0)
	r := bufio.NewReader(c.output)
	buf := make([]byte, 0, bytes.MinRead)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		nChunks++
		nBytes += int64(len(buf))
		c.content.Write(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		c.widget.SetText(c.content.String())
	}
	return
}

func (c *Client) Widget() *OutputWidget {
	return c.widget
}

func (c *Client) build(ev fsnotify.Event) error {
	builderCtx, cancel := context.WithCancel(c.builder.Ctx)
	done := make(chan struct{})
	defer func() {
		close(done)
	}()
	go func() {
		ticker := time.Tick(time.Second)
		for {
			select {
			case <-ticker:
				c.widget.SetBuildStatus(c.builder.State)
			case <-done:
				cancel()
				return
			}
		}
	}()
	err := c.builder.Rebuild(builderCtx)
	return err
}

func (c *Client) run(ev fsnotify.Event) error {
	// rerun the client
	return nil
}

func (c *Client) Notify(ev fsnotify.Event) {
	if err := c.build(ev); err != nil {
		return
	}
}

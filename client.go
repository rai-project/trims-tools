package micro

import (
	"bufio"
	"bytes"
	"context"
	"io"

	"github.com/fsnotify/fsnotify"
	tui "github.com/marcusolsson/tui-go"
)

type Client struct {
	ctx     context.Context
	path    string
	output  io.Reader
	content string
	builder *Builder
	watcher *watcher
	widget  *OutputWidget
}

func NewClient(ctx context.Context, ui tui.UI) *Client {
	ctx = context.WithValue(ctx, "name", "client")

	outputReader, outputWriter := io.Pipe()

	widget := NewOutputWidget("Client", ui)

	client := &Client{
		output:  outputReader,
		content: "",
		widget:  widget,
		builder: NewBuilder(ctx, outputWriter, outputWriter, Config.ClientPath, "make"),
	}
	defer func() {
		go client.listen()
	}()
	client.watcher = NewWatcher(
		client.Notify,
		Config.ClientPath,
		Config.ServerPath,
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
		c.content += string(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		c.widget.SetText(c.content)
	}
	return
}

func (c *Client) Widget() *OutputWidget {
	return c.widget
}

func (c *Client) Notify(ev fsnotify.Event) {
	err := c.builder.Rebuild()
	if err != nil {

	}
}

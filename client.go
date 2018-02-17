package micro

import (
	"github.com/dc0d/dirwatch"
	"github.com/drhodes/golorem"
	"github.com/fsnotify/fsnotify"
)

type Client struct {
	path       string
	content    string
	builder    *Builder
	watcher    *dirwatch.Watch
	widget     *OutputWidget
	isSelected bool
}

func NewClient() *Client {
	widget := NewOutputWidget("client")
	text := ""
	for ii := 0; ii < 100; ii++ {
		text += lorem.Sentence(89, 120) + "\n"
	}
	widget.SetText(text)
	return &Client{
		content: text,
		widget:  widget,
		builder: &Builder{},
	}
}

func (c *Client) Widget() *OutputWidget {
	return c.widget
}

func (c *Client) IsSelected(b bool) {
	c.isSelected = b
}

func (c *Client) notify(ev fsnotify.Event) {

}

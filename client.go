package micro

import (
	"github.com/dc0d/dirwatch"
	"github.com/drhodes/golorem"
	"github.com/fsnotify/fsnotify"
	tui "github.com/marcusolsson/tui-go"
)

type Client struct {
	path    string
	content *tui.Label
	builder *Builder
	watcher *dirwatch.Watch
}

func NewClient() *Client {
	content := tui.NewLabel("")
	text := ""
	for ii := 0; ii < 100; ii++ {
		text += lorem.Sentence(89, 120) + "\n"
	}
	content.SetText(text)
	return &Client{
		content: content,
		builder: &Builder{},
	}
}

func (c *Client) Widget() tui.Widget {
	content := tui.NewScrollArea(c.content)
	builder := c.builder.Widget()
	box := tui.NewHBox(
		content,
		builder,
	)
	box.SetBorder(true)
	box.SetTitle("client")

	box.SetSizePolicy(tui.Expanding, tui.Expanding)
	return box
}

func (c *Client) notify(ev fsnotify.Event) {

}

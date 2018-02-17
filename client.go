package micro

import (
	"time"

	"github.com/dc0d/dirwatch"
	"github.com/drhodes/golorem"
	"github.com/fsnotify/fsnotify"
	tui "github.com/marcusolsson/tui-go"
)

type Client struct {
	path       string
	content    string
	builder    *Builder
	watcher    *dirwatch.Watch
	widget     *OutputWidget
	isSelected bool
}

func NewClient(ui tui.UI) *Client {
	widget := NewOutputWidget("client", ui)
	go func() {
		ticker := time.Tick(time.Second)
		for {
			select {
			case <-ticker:
				text := ""
				for ii := 0; ii < 10; ii++ {
					text += lorem.Sentence(89, 120) + "\n"
				}
				widget.SetText(text)
			}
		}
	}()
	return &Client{
		content: "",
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

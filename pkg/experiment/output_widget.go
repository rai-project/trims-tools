package experiment

import (
	"fmt"
	"image"
	"strings"

	"github.com/marcusolsson/tui-go"
	"github.com/rai-project/micro18-tools/pkg/builder"
)

type OutputWidget struct {
	*tui.Box
	content       *tui.Label
	scrollArea    *tui.ScrollArea
	status        *tui.StatusBar
	contentText   string
	statusText    string
	scrollTopLeft image.Point

	ui tui.UI
}

func NewOutputWidget(title string, ui tui.UI) *OutputWidget {

	contentText := ""
	statusText := fmt.Sprintf("%s", "?")

	content := tui.NewLabel(contentText)
	content.SetWordWrap(true)

	scrollArea := tui.NewScrollArea(content)
	statusBar := tui.NewStatusBar(statusText)
	statusBar.SetPermanentText(title + " Build Status")
	box := tui.NewVBox(
		scrollArea,
		// HorizontalDivider,
		statusBar,
	)
	box.SetBorder(true)
	box.SetTitle(title)

	return &OutputWidget{
		content:       content,
		scrollArea:    scrollArea,
		status:        statusBar,
		Box:           box,
		contentText:   contentText,
		statusText:    statusText,
		scrollTopLeft: image.Point{0, 0},
		ui:            ui,
	}
}

func (w *OutputWidget) Scroll(dx, dy int) {
	w.scrollArea.Scroll(dx, dy)
	go w.ui.Update(func() {})
}

func (w *OutputWidget) SetBuildStatus(s builder.BuildState) {
	w.status.SetText(s.String())
	w.statusText = s.String()
}

func (w *OutputWidget) SetText(s string) {
	lineCount := len(strings.Split(s, "\n"))
	toScroll := lineCount - w.scrollArea.Size().Y + 2
	log.Infof("linecount = %d , w.scrollTopLeft.Y = %d , lineCount - w.scrollTopLeft.Y = %d", toScroll, w.scrollTopLeft.Y, toScroll-w.scrollTopLeft.Y)
	defer func() {
		w.scrollTopLeft.Y = toScroll
	}()
	w.contentText = s
	w.content.SetText(s)
	w.Scroll(w.scrollTopLeft.X, toScroll-w.scrollTopLeft.Y)
}

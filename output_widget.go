package micro

import (
	"fmt"
	"image"

	"github.com/k0kubun/pp"
	"github.com/marcusolsson/tui-go"
)

type OutputWidget struct {
	content     *tui.Label
	scrollArea  *tui.ScrollArea
	status      *tui.StatusBar
	box         *tui.Box
	contentText string
	statusText  string

	ui tui.UI
	tui.WidgetBase
}

var (
	divider = tui.NewLabel("-------------------")
)

func NewOutputWidget(title string, ui tui.UI) *OutputWidget {

	contentText := ""
	statusText := fmt.Sprintf("Build status : %s", "?")

	content := tui.NewLabel(contentText)
	content.SetWordWrap(true)

	scrollArea := tui.NewScrollArea(content)
	statusBar := tui.NewStatusBar(statusText)
	statusBox := tui.NewHBox(statusBar)
	statusBox.SetTitle("Status: ")
	statusBox.SetBorder(true)

	box := tui.NewVBox(
		scrollArea,
		divider,
		statusBox,
	)
	box.SetBorder(true)
	box.SetTitle(title)

	return &OutputWidget{
		content:     content,
		scrollArea:  scrollArea,
		status:      statusBar,
		box:         box,
		contentText: contentText,
		statusText:  statusText,
		ui:          ui,
	}
}

func (w *OutputWidget) Scroll(dx, dy int) {
	w.scrollArea.Scroll(dx, dy)
	go w.ui.Update(func() {})
}

func (w *OutputWidget) SetText(s string) {
	w.contentText = s
	w.content.SetText(s)
	go w.ui.Update(func() {})
}

func (w *OutputWidget) Draw(p *tui.Painter) {
	log.Info(pp.Sprint(w.box))
	w.box.Draw(p)
}

func (w *OutputWidget) SizeHint() image.Point {
	return w.box.SizeHint()
}

// MinSizeHint returns the size below which the widget cannot shrink.
func (w *OutputWidget) MinSizeHint() image.Point {
	return w.box.MinSizeHint()
}

// Size returns the current size of the widget.
func (w *OutputWidget) Size() image.Point {
	return w.box.Size()
}

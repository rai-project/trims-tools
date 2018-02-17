package micro

import (
	"fmt"
	"image"

	"github.com/marcusolsson/tui-go"
)

type OutputWidget struct {
	content    *tui.Label
	scrollArea *tui.ScrollArea
	status     *tui.StatusBar
	box        *tui.Box

	contentText string
	statusText  string

	tui.WidgetBase
}

var (
	divider = tui.NewLabel("-------------------")
)

func NewOutputWidget(title string) *OutputWidget {

	contentText := ""
	statusText := fmt.Sprintf("Build status : %s", "?")

	content := tui.NewLabel(contentText)
	scrollArea := tui.NewScrollArea(content)
	statusBar := tui.NewStatusBar("Status")
	statusBar.SetText(statusText)
	statusBox := tui.NewVBox(statusBar)
	statusBox.SetTitle("Status")
	statusBox.SetBorder(true)
	statusBar.SetText(statusText)

	box := tui.NewHBox(
		scrollArea,
		divider,
		statusBox,
	)
	box.SetBorder(true)
	box.SetTitle(title)

	box.SetSizePolicy(tui.Expanding, tui.Expanding)

	return &OutputWidget{
		content:     content,
		scrollArea:  scrollArea,
		status:      statusBar,
		box:         box,
		contentText: contentText,
		statusText:  statusText,
	}
}

func (w *OutputWidget) Scroll(dx, dy int) {
	w.scrollArea.Scroll(dx, dy)
}

func (w *OutputWidget) SetText(s string) {
	w.contentText = s
	w.content.SetText(s)
}

func (w *OutputWidget) Draw(p *tui.Painter) {
	w.box.Draw(p)
}

func (w *OutputWidget) SizeHint() image.Point {
	return image.Pt(len(w.contentText), 1)
}

package micro

import (
	"fmt"
	"os"

	"github.com/marcusolsson/tui-go"
)

type Term struct {
	ui     tui.UI
	theme  *tui.Theme
	server *Server
	client *Client
}

func NewTerm() *Term {
	tui.SetLogger(log)

	box := tui.NewHBox()
	box.SetTitle("micro18")
	box.SetBorder(true)

	ui, err := tui.New(box)
	if err != nil {
		panic(err)
	}

	theme := getTheme()
	ui.SetTheme(theme)

	server := NewServer(ui)
	client := NewClient(ui)

	serverWidget := server.Widget()
	clientWidget := client.Widget()

	_ = serverWidget
	_ = clientWidget

	box.Append(clientWidget)
	box.Append(clientWidget)

	tui.DefaultFocusChain.Set(serverWidget, clientWidget)

	ui.SetKeybinding("Left", func() {
		clientWidget.SetFocused(true)

		ui.SetKeybinding("Up", func() { clientWidget.Scroll(0, -1) })
		ui.SetKeybinding("Down", func() { clientWidget.Scroll(0, 1) })
		ui.SetKeybinding("Left", func() { clientWidget.Scroll(-1, 0) })
		ui.SetKeybinding("Right", func() { clientWidget.Scroll(1, 0) })
	})
	ui.SetKeybinding("Right", func() {
		// server.IsSelected()
	})

	ui.SetKeybinding("Esc", func() { quit(ui) })
	ui.SetKeybinding("q", func() { quit(ui) })
	ui.SetKeybinding("Ctrl+C", func() { quit(ui) })

	return &Term{
		ui:     ui,
		theme:  theme,
		server: server,
		client: client,
	}
}

func quit(ui tui.UI) {
	ui.Quit()
	fmt.Println("exit")
	os.Exit(0)
}

func getTheme() *tui.Theme {
	th := tui.DefaultTheme
	th.SetStyle("table.cell.selected", tui.Style{Bg: tui.ColorGreen, Fg: tui.ColorWhite})
	th.SetStyle("list.item", tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite})
	th.SetStyle("list.item.selected", tui.Style{Bg: tui.ColorGreen, Fg: tui.ColorWhite})
	return th
}

func (t *Term) Run() {
	if err := t.ui.Run(); err != nil {
		panic(err)
	}
}

var (
	VerticalDivider   *tui.Label
	HorizontalDivider *tui.Label
)

func init() {
	vline := ""
	hline := ""
	for ii := 0; ii < 10; ii++ {
		vline += VerticalLine.String()
		hline += HorizontalLine.String()
	}
	VerticalDivider = tui.NewLabel(vline)
	HorizontalDivider = tui.NewLabel(hline)
}

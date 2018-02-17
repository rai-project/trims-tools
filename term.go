package micro

import "github.com/marcusolsson/tui-go"

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

	box.Append(clientWidget)

	tui.DefaultFocusChain.Set(serverWidget, clientWidget)

	ui.SetKeybinding("Left", func() {
		client.IsSelected(true)

		ui.SetKeybinding("Up", func() { clientWidget.Scroll(0, -1) })
		ui.SetKeybinding("Down", func() { clientWidget.Scroll(0, 1) })
		ui.SetKeybinding("Left", func() { clientWidget.Scroll(-1, 0) })
		ui.SetKeybinding("Right", func() { clientWidget.Scroll(1, 0) })
	})
	ui.SetKeybinding("Right", func() {
		// server.IsSelected()
	})

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("q", func() { ui.Quit() })
	ui.SetKeybinding("Ctrl+C", func() { ui.Quit() })

	return &Term{
		ui:     ui,
		theme:  theme,
		server: server,
		client: client,
	}
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

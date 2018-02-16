package micro

import (
	"fmt"

	tui "github.com/marcusolsson/tui-go"
)

type BuildState int

const (
	BuildStateBuilding BuildState = 0
	BuildStateSuccess  BuildState = 1
	BuildStateFailed   BuildState = 2
)

type Builder struct {
	state BuildState
}

func (b *Builder) Widget() tui.Widget {
	state := map[BuildState]string{
		BuildStateFailed:   "∅",
		BuildStateBuilding: "…",
		BuildStateSuccess:  "✓",
	}

	label, ok := state[b.state]
	if !ok {
		label = "?"
	}

	statusBar := tui.NewStatusBar("")
	statusBox := tui.NewVBox(statusBar)
	statusBox.SetTitle("Status")
	statusBox.SetBorder(true)

	statusBar.SetText(fmt.Sprintf("Build status : %s", label))

	return statusBar
}

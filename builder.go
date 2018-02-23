package micro

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	tui "github.com/marcusolsson/tui-go"
	shellwords "github.com/mattn/go-shellwords"
)

type BuildState int

const (
	BuildStateBuilding BuildState = 0
	BuildStateSuccess  BuildState = 1
	BuildStateFailed   BuildState = 2
)

type Builder struct {
	ctx      context.Context
	state    BuildState
	baseDir  string
	buildCmd string
	stderr   io.Writer
	stdout   io.Writer
}

func NewBuilder(ctx context.Context, stderr io.Writer, stdout io.Writer, baseDir string, buildCmd string) *Builder {
	return &Builder{
		ctx:      context.WithValue(ctx, "name", "builder"),
		stderr:   stderr,
		stdout:   stdout,
		baseDir:  baseDir,
		buildCmd: buildCmd,
	}
}

func (b *Builder) Rebuild() error {
	b.state = BuildStateBuilding
	args, err := shellwords.Parse(b.buildCmd)
	if err != nil {
		log.WithError(err).WithField("cmd", b.buildCmd).Error("failed to parse shell command")
		return err
	}

	ctx, cancel := context.WithTimeout(b.ctx, time.Duration(Config.BuildTimeoutSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Path = b.baseDir
	if b.stdout != nil {
		cmd.Stdout = b.stdout
	} else {
		cmd.Stdout = os.Stdout
	}
	if b.stderr != nil {
		cmd.Stderr = b.stderr
	} else {
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		msg := fmt.Sprintf("failed to run %s within %s", b.buildCmd, b.baseDir)
		log.WithError(err).WithField("cmd", b.buildCmd).WithField("base", b.baseDir).Error("failed to run command")
		cmd.Stderr.Write([]byte(msg))
		b.state = BuildStateFailed
		return err
	}

	b.state = BuildStateSuccess
	return nil
}

func (b *Builder) Widget(ui tui.UI) tui.Widget {
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

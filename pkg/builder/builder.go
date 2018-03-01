package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	shellwords "github.com/mattn/go-shellwords"
)

type BuildState int

const (
	BuildStateBuilding BuildState = 0
	BuildStateSuccess  BuildState = 1
	BuildStateFailed   BuildState = 2
)

func (b BuildState) String() string {
	state := map[BuildState]string{
		BuildStateFailed:   "∅",
		BuildStateBuilding: "…",
		BuildStateSuccess:  "✓",
	}
	val, ok := state[b]
	if !ok {
		return "?"
	}
	return val
}

type Builder struct {
	Ctx      context.Context
	State    BuildState
	BaseDir  string
	BuildCmd string
	Stderr   io.Writer
	Stdout   io.Writer
	Cmd      *exec.Cmd
}

func (b *Builder) Rebuild(ctx context.Context) error {
	b.State = BuildStateBuilding
	args, err := shellwords.Parse(b.BuildCmd)
	if err != nil {
		log.WithError(err).WithField("cmd", b.BuildCmd).Error("failed to parse shell command")
		return err
	}
	if ctx == nil {
		ctx = b.Ctx
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(Config.BuildTimeoutSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Path = b.BaseDir
	if b.Stdout != nil {
		cmd.Stdout = b.Stdout
	} else {
		cmd.Stdout = os.Stdout
	}
	if b.Stderr != nil {
		cmd.Stderr = b.Stderr
	} else {
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		msg := fmt.Sprintf("failed to run %s within %s", b.BuildCmd, b.BaseDir)
		log.WithError(err).WithField("cmd", b.BuildCmd).WithField("base", b.BaseDir).Error("failed to run command")
		cmd.Stderr.Write([]byte(msg))
		b.State = BuildStateFailed
		return err
	}

	b.State = BuildStateSuccess
	return nil
}

package micro

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/github/hub/cmd"
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
	ctx      context.Context
	state    BuildState
	baseDir  string
	buildCmd string
	stderr   io.Writer
	stdout   io.Writer
	cmd      *cmd.Cmd
}

func (b *Builder) Rebuild(ctx context.Context) error {
	b.state = BuildStateBuilding
	args, err := shellwords.Parse(b.buildCmd)
	if err != nil {
		log.WithError(err).WithField("cmd", b.buildCmd).Error("failed to parse shell command")
		return err
	}
	if ctx == nil {
		ctx = b.ctx
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(Config.BuildTimeoutSeconds)*time.Second)
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

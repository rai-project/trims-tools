package config

import (
	"time"

	"github.com/rai-project/config"
	git "gopkg.in/src-d/go-git.v4"
)

type version struct {
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
}

var (
	Version version
)

func getVersion() {
	r, err := git.PlainOpen(Config.BaseSrcPath)
	if err != nil {
		panic(err)
	}

	h, err := r.Head()
	if err != nil {
		panic(err)
	}
	commit, err := r.CommitObject(h.Hash())
	if err != nil {
		panic(err)
	}

	wt, err := r.Worktree()
	if err != nil {
		panic(err)
	}

	st, err := wt.Status()
	if err != nil {
		panic(err)
	}

	state := "dirty"
	if st.IsClean() {
		state = "clean"
	}

	Version = version{
		BuildDate:  time.Now().String(),
		GitCommit:  h.Hash().String(),
		GitBranch:  h.Target().String(),
		GitState:   state,
		GitSummary: commit.Message,
	}
}

func init() {
	config.AfterInit(func() {
		Config.Wait()
		getVersion()
	})
}

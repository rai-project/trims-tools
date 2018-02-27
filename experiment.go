package micro

import (
	"time"

	"github.com/rai-project/config"
	"gopkg.in/mgo.v2/bson"
	git "gopkg.in/src-d/go-git.v4"
)

type Version struct {
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
}

var (
	RepositoryVersion Version
)

type Experiment struct {
	ID            bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt     time.Time     `json:"created_at"  bson:"created_at"`
	Hostname      string
	NetworkName   string
	FrameworkName string
	Version       Version

	Metadata map[string]string
}

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

	RepositoryVersion = Version{
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

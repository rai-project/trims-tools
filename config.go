package micro

import (
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/vipertags"
)

type microConfig struct {
	PollingInterval int           `json:"polling_interval" config:"micro18.polling_interval" default:100`
	Path            string        `json:"path" config:"micro18.path"`
	done            chan struct{} `json:"-" config:"-"`
}

var (
	Config = &microConfig{
		done: make(chan struct{}),
	}
)

func (microConfig) ConfigName() string {
	return "micro"
}

func (a *microConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

func (a *microConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
	if a.Path == "" {
		gopath := com.GetGOPATHs()[0]
		a.Path = filepath.Join(gopath, "src", "github.com", "rai-project", "mxnet-mirror")
	}
}

func (c microConfig) Wait() {
	<-c.done
}

func (c microConfig) String() string {
	return pp.Sprintln(c)
}

func (c microConfig) Debug() {
	log.Debug("micro Config = ", c)
}

func init() {
	config.Register(Config)
}

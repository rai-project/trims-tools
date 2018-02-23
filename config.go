package micro

import (
	"fmt"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/vipertags"
)

var (
	DefaultPath               = "Automatic"
	DefaultServerRelativePath = filepath.Join("bin", "uprd")
	DefaultClientRelativePath = filepath.Join("example", "image-classification", "predict-cpp")
)

type microConfig struct {
	BuildTimeoutSeconds int64         `json:"build_timeout" config:"micro18.build_timeout" default:600`
	PollingInterval     int           `json:"polling_interval" config:"micro18.polling_interval" default:100`
	BasePath            string        `json:"path" config:"micro18.path"`
	ServerRelativePath  string        `json:"server_relative_path" config:"micro18.server_relative_path"`
	ServerPath          string        `json:"server_path" config:"-"`
	ClientRelativePath  string        `json:"client_relative_path" config:"micro18.client_relative_path"`
	ClientPath          string        `json:"client_path" config:"-"`
	done                chan struct{} `json:"-" config:"-"`
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
	if a.BasePath == "" {
		if DefaultPath == "Automatic" {
			gopath := com.GetGOPATHs()[0]
			a.BasePath = filepath.Join(gopath, "src", "github.com", "rai-project", "mxnet-mirror")
		} else {
			a.BasePath = DefaultPath
		}
	}
	if !com.IsDir(a.BasePath) {
		panic(
			fmt.Sprintf("the directory %s does not exist. make sure that micro18.path is set correctly", a.BasePath))
	}
	if a.ServerRelativePath == "" {
		a.ServerRelativePath = DefaultServerRelativePath
	}
	if a.ClientRelativePath == "" {
		a.ClientRelativePath = DefaultClientRelativePath
	}
	a.ClientPath = filepath.Join(a.BasePath, a.ClientRelativePath)
	if !com.IsDir(a.ClientPath) {
		panic(
			fmt.Sprintf("the directory %s does not exist. make sure that micro18.path"+
				"and micro18.client_relative_path are set correctly", a.ClientPath,
			))
	}
	a.ServerPath = filepath.Join(a.BasePath, a.ServerRelativePath)
	if !com.IsDir(a.ServerPath) {
		panic(
			fmt.Sprintf("the directory %s does not exist. make sure that micro18.path"+
				"and micro18.server_relative_path are set correctly", a.ServerPath,
			))
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

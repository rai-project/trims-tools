package config

import (
	"os"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/k0kubun/pp"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rai-project/config"
	"github.com/rai-project/micro18-tools/pkg/utils"
	"github.com/rai-project/vipertags"
)

var (
	HomeDir, _                = homedir.Dir()
	DefaultSrcPath            = "Automatic"
	DefaultBasePath           = utils.GetEnvOr("UPR_BASE_DIR", filepath.Join(HomeDir, "carml", "data", "mxnet"))
	DefaultServerRelativePath = filepath.Join("bin", "uprd")
	DefaultServerBuildCmd     = "make"
	DefaultServerRunCmd       = filepath.Join("bin", "uprd")
	DefaultClientRelativePath = filepath.Join("example", "image-classification", "predict-cpp")
	DefaultClientBuildCmd     = "make"
	DefaultClientRunCmd       = "image-classification-predict"
)

type microConfig struct {
	BuildTimeoutSeconds int64         `json:"build_timeout" config:"micro18.build_timeout" default:600`
	PollingInterval     int           `json:"polling_interval" config:"micro18.polling_interval" default:100`
	BaseSrcPath         string        `json:"src_path" config:"micro18.src_path"`
	BasePath            string        `json:"base_path" config:"micro18.base_path"`
	ServerRelativePath  string        `json:"server_relative_path" config:"micro18.server_relative_path"`
	ServerPath          string        `json:"server_path" config:"-"`
	ServerBuildCmd      string        `json:"server_build_cmd" config:"micro18.server_build_cmd"`
	ServerRunCmd        string        `json:"server_run_cmd" config:"micro18.server_run_cmd"`
	ClientRelativePath  string        `json:"client_relative_path" config:"micro18.client_relative_path"`
	ClientPath          string        `json:"client_path" config:"-"`
	ClientBuildCmd      string        `json:"client_build_cmd" config:"micro18.client_build_cmd"`
	ClientRunCmd        string        `json:"client_run_cmd" config:"micro18.client_run_cmd"`
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
		a.BasePath = DefaultBasePath
		if !com.IsDir(a.BasePath) {
			err := os.MkdirAll(a.BasePath, os.ModePerm)
			if err != nil {
				log.WithError(err).WithField("base_path", a.BasePath).Panic("failed to create base path directory")
			}
		}
	}
	if a.BaseSrcPath == "" {
		if DefaultSrcPath == "Automatic" {
			gopath := com.GetGOPATHs()[0]
			a.BaseSrcPath = filepath.Join(gopath, "src", "github.com", "rai-project", "mxnet-mirror")
		} else {
			a.BaseSrcPath = DefaultSrcPath
		}
	}
	if !com.IsDir(a.BaseSrcPath) {
		log.Panicf("the directory %s does not exist. make sure that micro18.path is set correctly", a.BaseSrcPath)
	}
	if a.ServerRelativePath == "" {
		a.ServerRelativePath = DefaultServerRelativePath
	}
	if a.ClientRelativePath == "" {
		a.ClientRelativePath = DefaultClientRelativePath
	}
	a.ClientPath = filepath.Join(a.BaseSrcPath, a.ClientRelativePath)
	if !com.IsDir(a.ClientPath) {
		log.Panicf("the directory %s does not exist. make sure that micro18.path"+
			"and micro18.client_relative_path are set correctly", a.ClientPath,
		)
	}
	if a.ClientBuildCmd == "" {
		a.ClientBuildCmd = DefaultClientBuildCmd
	}
	if a.ClientRunCmd == "" {
		a.ClientRunCmd = DefaultClientRunCmd
	}
	a.ServerPath = filepath.Join(a.BaseSrcPath, a.ServerRelativePath)
	if !com.IsDir(a.ServerPath) {
		log.Panicf("the directory %s does not exist. make sure that micro18.path"+
			"and micro18.server_relative_path are set correctly", a.ServerPath,
		)
	}
	if a.ServerBuildCmd == "" {
		a.ServerBuildCmd = DefaultServerBuildCmd
	}
	if a.ServerRunCmd == "" {
		a.ServerRunCmd = DefaultServerRunCmd
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
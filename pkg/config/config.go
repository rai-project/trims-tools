package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/Unknwon/com"
	"github.com/k0kubun/pp"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rai-project/config"
	"github.com/rai-project/micro18-tools/pkg/utils"
	"github.com/rai-project/vipertags"
)

var (
	HomeDir, _                        = homedir.Dir()
	HostName, _                       = os.Hostname()
	DefaultVisibleDevices             = "0"
	DefaultSrcPath                    = "Automatic"
	DefaultBasePath                   = utils.GetEnvOr("UPR_BASE_DIR", filepath.Join(HomeDir, "carml", "data", "mxnet"))
	DefaultServerRelativePath         = ""
	DefaultServerBuildCmd             = "make"
	DefaultServerRunCmd               = filepath.Join("bin", "uprd")
	DefaultClientRelativePath         = filepath.Join("example", "image-classification", "predict-cpp")
	DefaultClientBuildCmd             = "make"
	DefaultClientRunCmd               = "./image-classification-predict"
	DefaultBaseBucketURL              = "http://s3.amazonaws.com/micro18profiles"
	DefaultUploadBucketName           = "traces"
	DefaultProfileOutputBaseDirectory = filepath.Join(HomeDir, "micro18_profiles")
	DefaultServerInfoPath             = filepath.Join(HomeDir, ".micro18_server.json")
)

type microConfig struct {
	VisibleDevices             string        `json:"cuda_visible_devices" yaml:"micro18.visible_devices" config:"micro18.visible_devices"`
	BuildTimeoutSeconds        int64         `json:"build_timeout" yaml:"micro18.build_timeout" config:"micro18.build_timeout" default:600`
	PollingInterval            int           `json:"polling_interval" yaml:"micro18.polling_interval" config:"micro18.polling_interval" default:100`
	BaseSrcPath                string        `json:"src_path" yaml:"micro18.src_path" config:"micro18.src_path"`
	BasePath                   string        `json:"base_path" yaml:"micro18.base_path" config:"micro18.base_path"`
	ServerRelativePath         string        `json:"server_relative_path" yaml:"micro18.server_relative_path" config:"micro18.server_relative_path"`
	ServerPath                 string        `json:"server_path" yaml:"-" config:"-"`
	ServerBuildCmd             string        `json:"server_build_cmd" yaml:"micro18.server_build_cmd" config:"micro18.server_build_cmd"`
	ServerRunCmd               string        `json:"server_run_cmd" yaml:"micro18.server_run_cmd" config:"micro18.server_run_cmd"`
	ClientRelativePath         string        `json:"client_relative_path" yaml:"micro18.client_relative_path" config:"micro18.client_relative_path"`
	ClientPath                 string        `json:"client_path" yaml:"-" config:"-"`
	ClientBuildCmd             string        `json:"client_build_cmd" yaml:"micro18.client_build_cmd" config:"micro18.client_build_cmd"`
	ClientRunCmd               string        `json:"client_run_cmd" yaml:"micro18.client_run_cmd" config:"micro18.client_run_cmd"`
	BaseBucketURL              string        `json:"base_bucket_url" yaml:"micro18.base_bucket_url" config:"micro18.base_bucket_url"`
	UploadBucketName           string        `json:"upload_bucket_name" yaml:"micro18.upload_bucket_name" config:"micro18.upload_bucket_name"`
	ProfileOutputBaseDirectory string        `json:"profile_output_base_directory" yaml:"micro18.profile_output_directory" config:"micro18.profile_output_directory"`
	ProfileOutputDirectory     string        `json:"profile_output_directory" yaml:"-" config:"-"`
	ExperimentDescription      string        `json:"experiment_description" yaml:"-" config:"-"`
	ServerInfoPath             string        `json:"server_info_path" config:"micro18.server_info_path"`
	UPREnabled                 bool          `json:"upr_enabled" config:"-"`
	done                       chan struct{} `json:"-" config:"-"`
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
	if a.ServerRelativePath == "" {
		a.ServerRelativePath = DefaultServerRelativePath
	}
	if a.ClientRelativePath == "" {
		a.ClientRelativePath = DefaultClientRelativePath
	}
	a.ClientPath = filepath.Join(a.BaseSrcPath, a.ClientRelativePath)
	if a.ClientBuildCmd == "" {
		a.ClientBuildCmd = DefaultClientBuildCmd
	}
	if a.ClientRunCmd == "" {
		a.ClientRunCmd = DefaultClientRunCmd
	}
	a.ServerPath = filepath.Join(a.BaseSrcPath, a.ServerRelativePath)
	if a.ServerBuildCmd == "" {
		a.ServerBuildCmd = DefaultServerBuildCmd
	}
	if a.ServerRunCmd == "" {
		a.ServerRunCmd = DefaultServerRunCmd
	}
	if a.UploadBucketName == "" {
		a.UploadBucketName = DefaultUploadBucketName
	}
	if a.BaseBucketURL == "" {
		a.BaseBucketURL = DefaultBaseBucketURL
	}
	if a.ProfileOutputBaseDirectory == "" {
		a.ProfileOutputBaseDirectory = DefaultProfileOutputBaseDirectory
	}
	if !com.IsDir(a.ProfileOutputBaseDirectory) {
		os.MkdirAll(a.ProfileOutputBaseDirectory, os.ModePerm)
	}
	a.ProfileOutputDirectory = filepath.Join(a.ProfileOutputBaseDirectory, HostName, time.Now().Format("2006-01-02-15"))
	if !com.IsDir(a.ProfileOutputDirectory) {
		os.MkdirAll(a.ProfileOutputDirectory, os.ModePerm)
	}
	if a.VisibleDevices == "" {
		a.VisibleDevices = DefaultVisibleDevices
	}
	if a.ServerInfoPath == "" {
		a.ServerInfoPath = DefaultServerInfoPath
	}
	a.UPREnabled = false
}

func (a microConfig) Verify() {
	if !com.IsDir(a.BaseSrcPath) {
		log.Panicf("the directory %s does not exist. make sure that micro18.path is set correctly", a.BaseSrcPath)
	}
	if !com.IsDir(a.ClientPath) {
		log.Panicf("the directory %s does not exist. make sure that micro18.path "+
			"and micro18.client_relative_path are set correctly", a.ClientPath,
		)
	}
	if !com.IsDir(a.ServerPath) {
		log.Panicf("the directory %s does not exist. make sure that micro18.path"+
			"and micro18.server_relative_path are set correctly", a.ServerPath,
		)
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

package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/Unknwon/com"
	"github.com/fatih/color"
	ps "github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/config"
	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/rai-project/micro18-tools/pkg/utils"
	"github.com/rai-project/uuid"
	"github.com/spf13/cast"
)

func IsRunning() bool {
	exe := config.Config.ServerRunCmd
	procs, err := ps.Processes()
	if err != nil {
		return false
	}
	for _, p := range procs {
		if p.Executable() == exe {
			return true
		}
		if p.Executable() == "uprd" {
			return true
		}
	}
	return false
}

func (s Server) IsRunning() bool {
	return IsRunning()
}

func (s *Server) Run() (*trace.Trace, error) {
	options := s.options

	cmdPath := filepath.Join(config.Config.ServerPath, config.Config.ServerRunCmd)
	if !com.IsFile(cmdPath) {
		return nil, errors.Errorf("the server command %s was not found in %s. make sure that the code compiled correctly",
			config.Config.ClientRunCmd, config.Config.ClientPath)
	}

	if ok, err := IsValidEvictionPolicy(options.evictionPolicy); !ok {
		return nil, err
	}

	fmt.Println(color.GreenString("✱ Running server and placing profile in " + config.Config.ProfileOutputDirectory))

	id := options.id
	if id == "" {
		id = uuid.NewV4()
	}
	profileFilePath := filepath.Join(config.Config.ProfileOutputDirectory, fmt.Sprintf("server_%s.json", id))
	env := map[string]string{
		"DATE":                               time.Now().Format(time.RFC3339Nano),
		"UPR_RUN_ID":                         id,
		"UPR_ENABLED":                        "true",
		"MXNET_ENGINE_TYPE":                  "ThreadedEngine",
		"UPR_PROFILE_TARGET":                 profileFilePath,
		"UPRD_EVICTION_POLICY":               fmt.Sprint(options.evictionPolicy),
		"UPRD_ESTIMATION_RATE":               fmt.Sprint(options.modelEstimationRate),
		"UPRD_MEMORY_PERCENTAGE":             fmt.Sprint(options.memoryPercentage),
		"UPR_GIT_SHA":                        config.Version.GitCommit,
		"UPR_GIT_BRANCH":                     config.Version.GitBranch,
		"UPR_GIT_Date":                       config.Version.BuildDate,
		"CUDA_VISIBLE_DEVICES":               config.Config.VisibleDevices,
		"UPR_BASE_DIR":                       config.Config.BasePath + "/",
		"MXNET_CUDNN_AUTOTUNE_DEFAULT":       "0",
		"UPRD_PERSIST_CPU":                   cast.ToString(options.persistCPU),
		"UPRD_ESTIMATE_WITH_INTERNAL_MEMORY": cast.ToString(options.estimateWithInternalMemory),
	}
	if options.debug {
		env["GLOG_logtostderr"] = "1"
		env["GLOG_v"] = "0"
		env["GLOG_stderrthreshold"] = "0"
	}
	if options.writeProfile {
		env["UPRD_WRITE_PROFILE"] = "true"
	} else {
		env["UPRD_WRITE_PROFILE"] = "false"
	}

	// log.WithField("server_path", config.Config.ServerPath).WithField("run_cmd", config.Config.ServerRunCmd).Debug("running server")
	ran, err := utils.ExecCmd(
		&s.cmd,
		config.Config.ServerPath,
		env,
		options.stdout,
		options.stderr,
		config.Config.ServerRunCmd,
	)
	if !ran {
		path := filepath.Join(config.Config.ServerPath, config.Config.ServerRunCmd)
		err = errors.Errorf("unable to run server cmd %s", path)
		log.WithError(err).Error("unable to run server")
		return nil, err
	}
	if err != nil {
		path := filepath.Join(config.Config.ServerPath, config.Config.ServerRunCmd)
		err = errors.Errorf("failed to run server cmd %s", path)
		log.WithField("cmd", config.Config.ServerRunCmd).WithError(err).Error("failed to run server")
		return nil, err
	}

	bts, err := ioutil.ReadFile(profileFilePath)
	if err != nil {
		err = errors.Wrapf(err, "unable to read profile file %s", profileFilePath)
		log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to read profile output")
		return nil, err
	}
	var trace trace.Trace
	if err := json.Unmarshal(bts, &trace); err != nil {
		err = errors.Wrapf(err, "unable to unmarshal profile file %s", profileFilePath)
		log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to unmarshal profile output")
		return nil, err
	}

	if options.uploadProfile {
		if err := trace.Upload(); err != nil {
			err = errors.Wrapf(err, "unable to upload profile file %s", profileFilePath)
			log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to upload profile output")
		}
	}

	return &trace, nil
}

func (s Server) Stop() error {
	if s.cmd == nil {
		return errors.New("command not started")
	}
	proc := s.cmd.Process
	return proc.Kill()
}

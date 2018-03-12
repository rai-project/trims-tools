package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/config"
	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/rai-project/micro18-tools/pkg/utils"
	"github.com/rai-project/uuid"
)

func Run(opts ...Option) (*trace.Trace, error) {
	options := WithOptions(opts...)

	if ok, err := IsValidEvictionPolicy(options.evictionPolicy); !ok {
		return nil, err
	}

	id := uuid.NewV4()
	profileFilePath := filepath.Join(config.Config.ProfileOutputDirectory, fmt.Sprintf("server_%s.json", id))
	env := map[string]string{
		"DATE":                   time.Now().Format(time.RFC3339Nano),
		"UPR_PROFILE_TARGET":     profileFilePath,
		"UPRD_EVICTION_POLICY":   fmt.Sprint(options.evictionPolicy),
		"UPRD_ESTIMATION_RATE":   fmt.Sprint(options.modelEstimationRate),
		"UPRD_MEMORY_PERCENTAGE": fmt.Sprint(options.memoryPercentage),
		"UPR_GIT_SHA":            config.Version.GitCommit,
		"UPR_GIT_BRANCH":         config.Version.GitBranch,
		"UPR_GIT_Date":           config.Version.BuildDate,
	}
	if options.debug {
		env["GLOG_logtostderr"] = "1"
		env["GLOG_v"] = "0"
		env["GLOG_stderrthreshold"] = "0"
	}

	ran, err := utils.ExecCmd(
		config.Config.ServerPath,
		env,
		os.Stdout,
		os.Stderr,
		config.Config.ServerRunCmd,
	)
	if !ran {
		err = errors.Errorf("unable to run server cmd %s", config.Config.ClientRunCmd)
		log.WithError(err).Error("unable to run server")
		return nil, err
	}
	if err != nil {
		err = errors.Errorf("failed to run server cmd %s", config.Config.ClientRunCmd)
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
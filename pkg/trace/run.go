package trace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rai-project/uuid"

	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/assets"
	"github.com/rai-project/micro18-tools/pkg/config"
)

func Run(opts ...Option) error {
	hostname, _ := os.Hostname()

	options := WithOptions(opts...)

	models := assets.Models
	if strings.ToLower(options.modelName) != "all" {
		models = assets.ModelManifests{}
		for _, m := range assets.Models {
			if strings.ToLower(m.MustCanonicalName()) == strings.ToLower(options.modelName) {
				models = assets.ModelManifests{m}
				break
			}
		}
		if len(models) == 0 {
			return errors.Errorf("the model %s was not found in the asset list", options.modelName)
		}
	}

	for _, model := range models {
		var combined *Trace
		for ii := 0; ii < options.iterationCount; ii++ {
			id := uuid.NewV4()
			profileFilePath := filepath.Join(config.Config.ProfileOutputDirectory, fmt.Sprintf("%s_%s.json", model.MustCanonicalName(), id))
			env := map[string]string{
				"DATE":                        time.Now().Format(time.RFC3339Nano),
				"UPR_MODEL_NAME":              model.MustCanonicalName(),
				"UPR_CLIENT":                  "1",
				"MXNET_CPU_PRIORITY_NTHREADS": "1",
				"OMP_NUM_THREADS":             "1",
				"MXNET_ENGINE_TYPE":           "NaiveEngine",
				"MXNET_GPU_WORKER_NTHREADS":   "1",
				"UPR_PROFILE_TARGET":          profileFilePath,
			}
			ran, err := execCmd(
				config.Config.ClientPath,
				env,
				os.Stdout,
				os.Stderr,
				config.Config.ClientRunCmd,
				fmt.Sprintf("%s_%d", hostname, ii),
			)
			if !ran {
				log.WithError(errors.Errorf("failed to run cmd %s", config.Config.ClientRunCmd)).Error("failed to run model")
				continue
			}
			if err != nil {
				log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to run model")
				continue
			}
			bts, err := ioutil.ReadFile(profileFilePath)
			if err != nil {
				err = errors.Wrapf(err, "unable to read profile file %s", profileFilePath)
				log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to read profile output")
				return err
			}
			var trace Trace
			if err := json.Unmarshal(bts, &trace); err != nil {
				err = errors.Wrapf(err, "unable to unmarshal profile file %s", profileFilePath)
				log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to unmarshal profile output")
				return err
			}
			trace.Iteration = int64(ii)
			if err := trace.Upload(); err != nil {
				err = errors.Wrapf(err, "unable to upload profile file %s", profileFilePath)
				log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to upload profile output")
			}
			if combined == nil {
				combined = &trace
				combined.ID = uuid.NewV4()
			} else {
				combined.Combine(trace)
			}
		}
	}
	return nil
}

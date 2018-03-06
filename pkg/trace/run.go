package trace

import (
	"fmt"
	"os"
	"strings"
	"time"

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
		for ii := 0; ii < options.iterationCount; ii++ {
			env := map[string]string{
				"DATE":                        time.Now().Format(time.RFC3339Nano),
				"UPR_MODEL_NAME":              model.MustCanonicalName(),
				"UPR_CLIENT":                  "1",
				"MXNET_CPU_PRIORITY_NTHREADS": "1",
				"OMP_NUM_THREADS":             "1",
				"MXNET_ENGINE_TYPE":           "NaiveEngine",
				"MXNET_GPU_WORKER_NTHREADS":   "1",
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
			}
			if err != nil {
				log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to run model")
			}
		}
	}
	return nil
}

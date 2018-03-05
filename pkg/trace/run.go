package trace

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/assets"
	"github.com/rai-project/micro18-tools/pkg/config"
)

func Run(ctx context.Context, ii int) error {
	models := assets.Models
	hostname, _ := os.Hostname()
	ran, err := execCmd(
		config.Config.ClientPath,
		map[string]string{
			"DATE":                        time.Now().Format(time.RFC3339Nano),
			"UPR_CLIENT":                  "1",
			"MXNET_CPU_PRIORITY_NTHREADS": "1",
			"OMP_NUM_THREADS":             "1",
			"MXNET_ENGINE_TYPE":           "NaiveEngine",
			"MXNET_GPU_WORKER_NTHREADS":   "1",
		},
		os.Stdout,
		os.Stderr,
		config.Config.ClientRunCmd,
		fmt.Sprintf("%s_%d", hostname, ii),
	)
	if !ran {
		return errors.Errorf("failed to run cmd %s", config.Config.ClientRunCmd)
	}
	return err
}

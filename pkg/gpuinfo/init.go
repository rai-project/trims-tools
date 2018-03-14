package gpuinfo

import (
	"runtime"

	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	log         = logger.New().WithField("pkg", "micro/gpuinfo")
	IsSupported = runtime.GOOS == "linux"
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "micro/gpuinfo")
	})
}

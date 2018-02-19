package nvprof

import (
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	log = logger.New().WithField("pkg", "micro/nvprof")
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "micro/nvprof")
	})
}

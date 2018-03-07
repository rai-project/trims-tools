package workload

import (
	"github.com/dc0d/dirwatch"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	log = logger.New().WithField("pkg", "micro/workload")
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "micro/workload")
		dirwatch.SetLogger(log.Debug)
	})
}

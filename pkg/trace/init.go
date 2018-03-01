package trace

import (
	"github.com/dc0d/dirwatch"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	log = logger.New().WithField("pkg", "micro/trace")
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "micro/trace")
		dirwatch.SetLogger(log)
	})
}

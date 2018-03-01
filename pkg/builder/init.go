package builder

import (
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	mconfig "github.com/rai-project/micro18-tools/pkg/config"
)

var (
	log    = logger.New().WithField("pkg", "micro/builder")
	Config = mconfig.Config
)

func init() {
	config.AfterInit(func() {
		Config.Wait()
		log = logger.New().WithField("pkg", "micro/builder")
	})
}

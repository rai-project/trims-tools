package server

import (
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	log = logger.New().WithField("pkg", "micro/server")
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "micro/server")
	})
}

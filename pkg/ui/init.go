package ui

import (
	"os"

	"github.com/dc0d/dirwatch"

	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	log = logger.New().WithField("pkg", "micro/term")
)

func init() {
	config.AfterInit(func() {
		lg := logger.New()
		f, _ := os.Create("debug.log")
		lg.Out = f
		log = lg.WithField("pkg", "micro/term")
		dirwatch.SetLogger(log.Info)
	})
}

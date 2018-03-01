// +build ignore

// #!/usr/bin/env gorun

package main

import (
	"context"

	"github.com/fatih/color"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	ui "github.com/rai-project/micro18-tools/pkg/ui"
)

func main() {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		logger.Error(e)
	}()

	color.NoColor = false
	opts := []config.Option{
		config.AppName("carml"),
		config.ColorMode(true),
	}

	config.Init(opts...)

	ui := ui.NewTerm(context.Background())
	ui.Run()
}

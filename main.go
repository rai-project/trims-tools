// +build ignore

package main

import (
	"github.com/fatih/color"
	"github.com/rai-project/config"
	micro "github.com/rai-project/micro18-tools"
)

func main() {

	color.NoColor = false
	opts := []config.Option{
		config.AppName("carml"),
		config.ColorMode(true),
	}

	config.Init(opts...)

	ui := micro.NewTerm()
	ui.Run()
}

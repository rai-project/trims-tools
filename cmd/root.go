package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/fatih/color"
	"github.com/k0kubun/pp"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	IsDebug   bool = true
	IsVerbose bool = true
	CfgFile   string
	log       *logrus.Entry = logrus.New().WithField("pkg", "micro/cmd")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "micro",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is $HOME/.carml_config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	log.Level = logrus.DebugLevel
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "micro/cmd")
	})

	color.NoColor = false
	opts := []config.Option{
		config.AppName("carml"),
		config.ColorMode(true),
		config.DebugMode(IsDebug),
		config.VerboseMode(IsVerbose),
	}
	if IsDebug || IsVerbose {
		pp.WithLineInfo = true
	}
	if c, err := homedir.Expand(CfgFile); err == nil {
		CfgFile = c
	}
	if config.IsValidRemotePrefix(CfgFile) {
		opts = append(opts, config.ConfigRemotePath(CfgFile))
	} else if com.IsFile(CfgFile) {
		if c, err := filepath.Abs(CfgFile); err == nil {
			CfgFile = c
		}
		opts = append(opts, config.ConfigFileAbsolutePath(CfgFile))
	}

	config.Init(opts...)

}

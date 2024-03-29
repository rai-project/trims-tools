package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Unknwon/com"
	"github.com/fatih/color"
	"github.com/k0kubun/pp"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	mconfig "github.com/rai-project/micro18-tools/pkg/config"
	"github.com/rai-project/micro18-tools/pkg/gpuinfo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	IsDebug                bool = true
	IsVerbose              bool = true
	profileOutputOverwrite bool
	CfgFile                string
	monitorMemory          bool
	memoryInfo             *gpuinfo.System
	visibleDevices         string
	profileOutput          string
	experimentDescription  string
	log                    *logrus.Entry = logrus.New().WithField("pkg", "micro/cmd")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "micro",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		mconfig.Config.VisibleDevices = visibleDevices
		os.Setenv("CUDA_VISIBLE_DEVICES", visibleDevices)
		if monitorMemory && gpuinfo.IsSupported {
			info, err := gpuinfo.New()
			if err != nil {
				pp.Println(err)
				log.WithError(err).Error("failed to create gpu memory info object")
				return nil
			}
			memoryInfo = info
			memoryInfo.Start(50 * time.Millisecond)
		}
		if profileOutput != "" {
			mconfig.Config.ProfileOutputDirectory = filepath.Join(mconfig.Config.ProfileOutputBaseDirectory, mconfig.HostName, profileOutput)
			if profileOutputOverwrite && com.IsDir(mconfig.Config.ProfileOutputDirectory) {
				os.RemoveAll(mconfig.Config.ProfileOutputDirectory)
			}
			if !com.IsDir(mconfig.Config.ProfileOutputDirectory) {
				os.MkdirAll(mconfig.Config.ProfileOutputDirectory, os.ModePerm)
			}
		}
		mconfig.Config.ExperimentDescription = experimentDescription
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if memoryInfo != nil {
			memoryInfo.Stop()
			memoryInfo.Print()
		}
		return nil
	},
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
	rootCmd.PersistentFlags().BoolVar(&profileOutputOverwrite, "profile_output_overwrite", false, "delete output directory for the profiles if it exists")
	rootCmd.PersistentFlags().StringVar(&profileOutput, "profile_output", "", "output directory for the profiles")
	rootCmd.PersistentFlags().StringVar(&experimentDescription, "experiment_description", "", "description of the experiement run")
	rootCmd.PersistentFlags().BoolVar(&monitorMemory, "monitor_memory", false && gpuinfo.IsSupported, "monitors the memory during evaluation and prints the memory information at the end")
	rootCmd.PersistentFlags().StringVar(&visibleDevices, "visible_devices", "0", "comma seperated list of devices visible to both the server and client. This controls the CUDA_VISIBLE_DEVICES variable")

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
	if com.IsFile(CfgFile) {
		if c, err := filepath.Abs(CfgFile); err == nil {
			CfgFile = c
		}
		opts = append(opts, config.ConfigFileAbsolutePath(CfgFile))
	}

	config.Init(opts...)

}

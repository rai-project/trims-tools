package cmd

import (
	"context"

	"github.com/rai-project/micro18-tools/pkg/client"
	"github.com/spf13/cobra"
)

var (
	runClientOriginal                    bool
	runClientNTimes                      int
	runClientEager                       bool
	runClientEagerAsync                  bool
	runClientPostprocess                 bool
	runClientDebug                       bool
	runClientModelDistribution           string
	runClientModelDistributionParameters string
	runClientConcurrentCount             int
)

// runCmd represents the run command
var clientRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the client command and produce profile files",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client := client.New(
			client.Context(ctx),
			client.OriginalMode(runClientOriginal),
			client.DebugMode(runClientDebug),
			client.PostProcess(runClientPostprocess),
			client.IterationCount(runClientNTimes),
			client.EagerInitialize(runClientEager),
			client.EagerInitializeAsync(runClientEagerAsync),
			client.ModelDistribution(runClientModelDistribution, runClientModelDistributionParameters),
		)
		_, err := client.Run()
		return err
	},
}

func init() {
	clientCmd.AddCommand(clientRunCmd)
	clientRunCmd.Flags().BoolVar(&runClientOriginal, "original", false, "Run an unmodified version of the inference (without persistent storage)")
	clientRunCmd.Flags().IntVarP(&runClientNTimes, "iterations", "n", 1, "Number of iterations to run the client")
	clientRunCmd.Flags().StringVar(&runClientModelDistribution, "distribution", "uniform", "Distribution for selecting models while running client")
	clientRunCmd.Flags().StringVar(&runClientModelDistributionParameters, "distribution_params", "", "Distribution parameters for selecting models while running client")
	clientRunCmd.Flags().IntVar(&runClientConcurrentCount, "concurrent", 1, " Number of clients to run concurrently")
	clientRunCmd.Flags().BoolVarP(&runClientDebug, "debug", "d", false, "Print debug messages from the client")
	clientRunCmd.Flags().BoolVar(&runClientPostprocess, "postprocess", true, "whether to postprocess the client output as part of the run")
	clientRunCmd.Flags().BoolVar(&runClientEager, "eager", true, "eagerly initialize the client")
	clientRunCmd.Flags().BoolVar(&runClientEagerAsync, "eager_async", false, "eagerly initialize the client but make it asynchronous")
}

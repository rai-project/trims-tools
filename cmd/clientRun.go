package cmd

import (
	"context"

	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/spf13/cobra"
)

var (
	runClientNTimes      int
	runClientEager       bool
	runClientEagerAsync  bool
	runClientPostprocess bool
)

// runCmd represents the run command
var clientRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the client command and produce profile files",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		trace.Run(
			trace.Context(ctx),
			trace.PostProcess(runClientPostprocess),
			trace.IterationCount(runClientNTimes),
			trace.EagerInitialize(runClientEager),
			trace.EagerInitializeAsync(runClientEagerAsync),
		)
	},
}

func init() {
	clientCmd.AddCommand(clientRunCmd)
	clientRunCmd.Flags().IntVarP(&runClientNTimes, "iterations", "n", 1, "Number of iterations to run the client")
	clientRunCmd.Flags().BoolVar(&runClientPostprocess, "postprocess", true, "whether to postprocess the client output as part of the run")
	clientRunCmd.Flags().BoolVar(&runClientEager, "eager", true, "eagerly initialize the client")
	clientRunCmd.Flags().BoolVar(&runClientEagerAsync, "eager_async", false, "eagerly initialize the client but make it asynchronous")
}

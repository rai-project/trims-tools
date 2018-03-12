package cmd

import (
	"context"

	"github.com/rai-project/micro18-tools/pkg/client"
	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/spf13/cobra"
)

var (
	runClientOriginal                    bool
	runClientProfileIO                   bool
	runClientNTimes                      int
	runClientEager                       bool
	runClientEagerAsync                  bool
	runClientPostprocess                 bool
	runClientDebug                       bool
	runClientModelDistribution           string
	runClientModelDistributionParameters string
	runClientConcurrentCount             int
	runClientModelIterations             int
)

func makeClientRun(ctx context.Context, extraOpts ...client.Option) *client.Client {
	opts := append(
		[]client.Option{
			client.Context(ctx),
			client.OriginalMode(runClientOriginal),
			client.ProfileIO(runClientProfileIO),
			client.DebugMode(runClientDebug),
			client.PostProcess(runClientPostprocess),
			client.IterationCount(runClientNTimes),
			client.EagerInitialize(runClientEager),
			client.EagerInitializeAsync(runClientEagerAsync),
			client.ConcurrentRunCount(runClientConcurrentCount),
			client.ModelIterationCount(runClientModelIterations),
			client.ModelDistribution(runClientModelDistribution, runClientModelDistributionParameters),
		},
		extraOpts...,
	)
	return client.New(opts...)
}

var clientRunCompare = &cobra.Command{
	Use:     "run-compare",
	Aliases: []string{"compare"},
	Short:   "Run the client command and produce profile files",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		orig := makeClientRun(ctx, client.OriginalMode(true))
		origTrace, err := orig.Run()

		mod := makeClientRun(ctx, client.OriginalMode(false))
		modTrace, err := mod.Run()

		firstTrace := *modTrace[0]
		restTraces := []trace.Trace{}
		for _, tr := range append(modTrace[1:], origTrace...) {
			restTraces = append(restTraces, *tr)
		}
		trace.Combine(firstTrace, restTraces...)

		return err
	},
}

// runCmd represents the run command
var clientRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the client command and produce profile files",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client := makeClientRun(ctx)
		_, err := client.Run()
		return err
	},
}

func init() {
	runCmds := []*cobra.Command{
		clientRunCmd,
		clientRunCompare,
	}
	for _, cmd := range runCmds {
		clientCmd.AddCommand(cmd)
		cmd.Flags().BoolVar(&runClientProfileIO, "profileio", true, "Profile I/O model read (this only makes sense when evaluating the original mxnet implementation)")
		cmd.Flags().IntVar(&runClientModelIterations, "model_iterations", -1, "Number of iterations to run each model")
		cmd.Flags().IntVarP(&runClientNTimes, "iterations", "n", 1, "Number of iterations to run the client")
		cmd.Flags().StringVar(&runClientModelDistribution, "distribution", "none", "Distribution for selecting models while running client")
		cmd.Flags().StringVar(&runClientModelDistributionParameters, "distribution_params", "", "Distribution parameters for selecting models while running client")
		cmd.Flags().IntVar(&runClientConcurrentCount, "concurrent", 1, " Number of clients to run concurrently")
		cmd.Flags().BoolVarP(&runClientDebug, "debug", "d", false, "Print debug messages from the client")
		cmd.Flags().BoolVar(&runClientPostprocess, "postprocess", true, "whether to postprocess the client output as part of the run")
		cmd.Flags().BoolVar(&runClientEager, "eager", true, "eagerly initialize the client")
		cmd.Flags().BoolVar(&runClientEagerAsync, "eager_async", false, "eagerly initialize the client but make it asynchronous")
	}
	clientRunCmd.Flags().BoolVar(&runClientOriginal, "original", false, "Run an unmodified version of the inference (without persistent storage)")
}

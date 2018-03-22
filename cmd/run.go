package cmd

import (
	"context"
	"time"

	mconfig "github.com/rai-project/micro18-tools/pkg/config"
	"github.com/spf13/cobra"
)

var (
	runModels                      string
	runCompare                     bool
	runOriginal                    bool
	runProfileIO                   bool
	runNTimes                      int
	runEager                       bool
	runEagerAsync                  bool
	runPostprocess                 bool
	runDebug                       bool
	runModelDistribution           string
	runModelDistributionParameters string
	runConcurrentCount             int
	runModelIterations             int
	runProfileMemory               bool
	runUploadTraces                bool
	runCombinedAll                 bool
	runEvictionPolicy              string
	runModelEstimationRate         float32
	runMemoryPercentage            float32
	runPersistCPU                  bool
	runWriteProfile                bool
	runLargeModels                 bool
	runEstimateWithInternalMemory  bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		runClientModels = runModels
		runClientDebug = runDebug
		runClientOriginal = runOriginal
		runClientProfileIO = runProfileIO
		runClientNTimes = runNTimes
		runClientEager = runEager
		runClientEagerAsync = runEagerAsync
		runClientPostprocess = runPostprocess
		runClientModelDistribution = runModelDistribution
		runClientModelDistributionParameters = runModelDistributionParameters
		runClientConcurrentCount = runConcurrentCount
		runClientModelIterations = runModelIterations
		runClientProfileMemory = runProfileMemory
		runClientProfileMemory = runProfileMemory
		runClientCombinedAll = runCombinedAll
		runClientLargeModels = runLargeModels

		runServerDebug = runDebug
		runServerEvictionPolicy = runEvictionPolicy
		runServerModelEstimationRate = runModelEstimationRate
		runServerMemoryPercentage = runMemoryPercentage
		runServerPersistCPU = runPersistCPU
		runServerWriteProfile = runWriteProfile
		runServerEstimateWithInternalMemory = runEstimateWithInternalMemory

		if err := serverRunCmd.PreRunE(cmd, args); err != nil {
			return err
		}
		if runOriginal {
			mconfig.Config.UPREnabled = false
		} else {
			mconfig.Config.UPREnabled = true
		}
		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return serverRunCmd.PostRunE(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
				return
			}
		}()

		ctx := context.Background()
		server := makeServerRun(ctx)
		defer server.Stop()

		go func() {
			_, err := server.Run()
			if err != nil {
				panic(err)
			}
		}()

		time.Sleep(time.Second)

		if runCompare {
			err = clientCompare(ctx)
			return
		}

		client := makeClientRun(ctx)
		_, err = client.Run(!runCompare)
		return
	},
}

func init() {
	runCmd.Flags().BoolVar(&runCompare, "compare", false, "Compare with and without proposal")
	runCmd.Flags().BoolVarP(&runDebug, "debug", "d", false, "Print debug messages from the client")
	runCmd.Flags().BoolVar(&runWriteProfile, "profile", false, "Write the server profile file")
	runCmd.Flags().StringVar(&runEvictionPolicy, "eviction", "lru", "Eviction policy used by the server")
	runCmd.Flags().Float32Var(&runModelEstimationRate, "model_estimation_rate", 1.0, "File size multiplier used to determine how much memory would be used by a model")
	runCmd.Flags().Float32Var(&runMemoryPercentage, "memory_percentage", 0.8, "Percentage of GPU memory that can be used to persist models")
	runCmd.Flags().BoolVar(&runPersistCPU, "persist_cpu", true, "Persist memory on CPU to avoid rereading the data from disk after eviction")
	runCmd.Flags().StringVar(&runModels, "models", "all", "List of models to use (comma seperated)")
	runCmd.Flags().BoolVar(&runProfileIO, "profile_io", true, "Profile I/O model read (this only makes sense when evaluating the original mxnet implementation)")
	runCmd.Flags().IntVar(&runModelIterations, "model_iterations", -1, "Maximum number of iterations to run each model")
	runCmd.Flags().IntVarP(&runNTimes, "iterations", "n", 1, "Number of iterations to run the client")
	runCmd.Flags().StringVar(&runModelDistribution, "distribution", "none", "Distribution for selecting models while running client")
	runCmd.Flags().StringVar(&runModelDistributionParameters, "distribution_params", "", "Distribution parameters for selecting models while running client")
	runCmd.Flags().IntVar(&runConcurrentCount, "concurrent", 1, " Number of clients to run concurrently")
	runCmd.Flags().BoolVar(&runPostprocess, "postprocess", true, "whether to postprocess the client output as part of the run")
	runCmd.Flags().BoolVar(&runEager, "eager", true, "eagerly initialize the client")
	runCmd.Flags().BoolVar(&runEagerAsync, "eager_async", false, "eagerly initialize the client but make it asynchronous")
	runCmd.Flags().BoolVar(&runProfileMemory, "profile_memory", true, "track the cudaMalloc and cudaFree calls")
	runCmd.Flags().BoolVar(&runUploadTraces, "trace_upload", false, "upload the traces to AWS S3 once complete")
	runCmd.Flags().BoolVar(&runOriginal, "original", false, "Run an unmodified version of the inference (without persistent storage)")
	runCmd.Flags().BoolVar(&runCombinedAll, "combined_all", true, "Combine all results into a single trace")
	runCmd.Flags().BoolVar(&runLargeModels, "large_models", false, "run the large models")
	runCmd.Flags().BoolVar(&runEstimateWithInternalMemory, "estimate_with_internal_memory", true, "Use internal memory information when estimating model size")

	rootCmd.AddCommand(runCmd)
}

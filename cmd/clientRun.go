package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/assets"
	"github.com/rai-project/micro18-tools/pkg/client"
	mconfig "github.com/rai-project/micro18-tools/pkg/config"
	"github.com/rai-project/micro18-tools/pkg/server"
	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/rai-project/micro18-tools/pkg/utils"
	"github.com/rai-project/uuid"
	"github.com/spf13/cobra"
)

var (
	runClientModels                      string
	runClientPercentageModels            float64
	runClientOriginal                    bool
	runClientCompareOriginal             bool
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
	runClientProfileMemory               bool
	runClientUploadTraces                bool
	runClientCombinedAll                 bool
	runClientLargeModels                 bool
	runClientSimulateRun                 bool
	runClientCompareCombineTraces        bool
)

func makeClientRun(ctx context.Context, extraOpts ...client.Option) *client.Client {
	opts := append(
		[]client.Option{
			client.Context(ctx),
			client.ModelName(runClientModels),
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
			client.ProfileMemory(runClientProfileMemory),
			client.UploadProfile(runClientUploadTraces),
			client.LargeModels(runClientLargeModels),
			client.SimulateRun(runClientSimulateRun),
			client.PercentageModels(runClientPercentageModels),
		},
		extraOpts...,
	)
	return client.New(opts...)
}

func clientCompare(ctx context.Context) error {
	combinedTraces := map[string][]*trace.Trace{}

	var models assets.ModelManifests
	var err error
	if runClientLargeModels {
		models, err = assets.FilterLargeModels(runClientModels)
	} else {
		models, err = assets.FilterModels(runClientModels)
	}
	if err != nil {
		return err
	}

	progressCount := len(models)
	if runClientCompareOriginal {
		progressCount *= 2
	}
	progress := utils.NewProgress("comparing models", progressCount)
	defer progress.FinishPrint("finished comparing models")

	originalTracesMap := map[string][]*trace.Trace{}
	modTracesMap := map[string][]*trace.Trace{}

	if runClientCompareOriginal {
		for _, model := range models {
			progress.Increment()
			//println("running ", model.MustCanonicalName(), " in original mode")
			orig := makeClientRun(ctx, client.OriginalMode(true), client.ModelName(model.MustCanonicalName()))
			origTraces, err := orig.Run(false)
			if err != nil {
				log.WithError(err).Error("failed to run client with upr enabled")
				continue
			}
			originalTracesMap[model.MustCanonicalName()] = origTraces
		}
	}

	for _, model := range models {
		progress.Increment()
		//println("running ", model.MustCanonicalName(), " in upr mode")
		mod := makeClientRun(ctx, client.OriginalMode(false), client.ModelName(model.MustCanonicalName()))
		modTraces, err := mod.Run(false)
		if err != nil {
			log.WithError(err).Error("failed to run client with upr disabled")
			continue
		}
		modTracesMap[model.MustCanonicalName()] = modTraces
	}

	for _, model := range models {
		name := model.MustCanonicalName()
		origTraces, ok := originalTracesMap[name]
		if !ok {
			continue
		}
		modTraces, ok := modTracesMap[name]
		if !ok {
			continue
		}
		if len(modTraces) == 0 {
			log.WithField("model_name", model.MustCanonicalName()).Error("no traces captured")
			continue
		}

		if runClientCombinedAll {
			if _, ok := combinedTraces["all"]; !ok {
				combinedTraces["all"] = []*trace.Trace{}
			}
			combinedTraces["all"] = append(combinedTraces["all"], origTraces...)
			combinedTraces["all"] = append(combinedTraces["all"], modTraces...)
		} else {
			name := model.MustCanonicalName()
			if _, ok := combinedTraces[name]; !ok {
				combinedTraces[name] = []*trace.Trace{}
			}
			combinedTraces[name] = append(combinedTraces[name], origTraces...)
			combinedTraces[name] = append(combinedTraces[name], modTraces...)
		}
	}
	if runClientCompareCombineTraces {
		progress.Prefix("combining traces")
		for name, traces := range combinedTraces {
			if len(traces) == 0 {
				continue
			}

			firstTrace := *traces[0]
			if tr, err := firstTrace.Adjust(); err == nil {
				firstTrace = tr
			}

			restTraces := []trace.Trace{}
			for _, tr := range traces[1:] {
				if tr == nil {
					continue
				}
				if adjustedTrace, err := tr.Adjust(); err == nil {
					restTraces = append(restTraces, adjustedTrace)
					continue
				}
				restTraces = append(restTraces, *tr)
			}
			combined := trace.Combine(firstTrace, restTraces...)

			if combined != nil {
				id := uuid.NewV4()
				path := filepath.Join(mconfig.Config.ProfileOutputDirectory, "compared-"+name+"-"+id+".json")
				bts, err := json.Marshal(combined)
				if err == nil {
					ioutil.WriteFile(path, bts, 0644)
				}
			}
		}
	}

	return err
}

var clientRunCompare = &cobra.Command{
	Use:     "run-compare",
	Aliases: []string{"compare"},
	Short:   "Run the client command and produce profile files",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !runClientOriginal && !server.IsRunning() {
			return errors.New("the uprd server is not running. make sure you've started the server before starting the client")
		}
		if runClientOriginal {
			mconfig.Config.UPREnabled = false
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		return clientCompare(ctx)
	},
}

// runCmd represents the run command
var clientRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the client command and produce profile files",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !runClientOriginal && !server.IsRunning() {
			return errors.New("the uprd server is not running. make sure you've started the server before starting the client")
		}
		if runClientOriginal {
			mconfig.Config.UPREnabled = false
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client := makeClientRun(ctx)
		_, err := client.Run(true)
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
		cmd.Flags().StringVar(&runClientModels, "models", "all", "List of models to use (comma seperated)")
		cmd.Flags().Float64Var(&runClientPercentageModels, "models_percentage", 1.0, "percentage of models to run")
		cmd.Flags().BoolVar(&runClientProfileIO, "profile_io", true, "Profile I/O model read (this only makes sense when evaluating the original mxnet implementation)")
		cmd.Flags().IntVar(&runClientModelIterations, "model_iterations", -1, "Maximum number of iterations to run each model")
		cmd.Flags().IntVarP(&runClientNTimes, "iterations", "n", 1, "Number of iterations to run the client")
		cmd.Flags().StringVar(&runClientModelDistribution, "distribution", "none", "Distribution for selecting models while running client")
		cmd.Flags().StringVar(&runClientModelDistributionParameters, "distribution_params", "", "Distribution parameters for selecting models while running client")
		cmd.Flags().IntVar(&runClientConcurrentCount, "concurrent", 1, " Number of clients to run concurrently")
		cmd.Flags().BoolVarP(&runClientDebug, "debug", "d", false, "Print debug messages from the client")
		cmd.Flags().BoolVar(&runClientPostprocess, "postprocess", true, "whether to postprocess the client output as part of the run")
		cmd.Flags().BoolVar(&runClientEager, "eager", true, "eagerly initialize the client")
		cmd.Flags().BoolVar(&runClientEagerAsync, "eager_async", false, "eagerly initialize the client but make it asynchronous")
		cmd.Flags().BoolVar(&runClientProfileMemory, "profile_memory", true, "track the cudaMalloc and cudaFree calls")
		cmd.Flags().BoolVar(&runClientUploadTraces, "trace_upload", false, "upload the traces to AWS S3 once complete")
		cmd.Flags().BoolVar(&runClientLargeModels, "large_models", false, "run the large models")
		cmd.Flags().BoolVar(&runClientSimulateRun, "simulate_run", false, "do not run the command, just simulate the process of running it")
	}
	clientRunCmd.Flags().BoolVar(&runClientOriginal, "original", false, "Run an unmodified version of the inference (without persistent storage)")
	clientRunCompare.Flags().BoolVar(&runClientCombinedAll, "combined_all", false, "Combine all results into a single trace")
	clientRunCompare.Flags().BoolVar(&runClientCompareOriginal, "run_original", true, "Run the original code when comparing")
	clientRunCompare.Flags().BoolVar(&runClientCompareCombineTraces, "combine", false, "Combine the comparison traces into one timeline")
}

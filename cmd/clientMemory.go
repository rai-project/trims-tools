package cmd

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/rai-project/micro18-tools/pkg/client"
	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/spf13/cobra"
)

var (
	runClientMemoryOutputFormat string
)

var clientRunMemoryCmd = &cobra.Command{
	Use:     "memory_usage",
	Aliases: []string{"memory"},
	Short:   "Run the client command and print out how much memory is used by internal layers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		opts :=
			[]client.Option{
				client.Context(ctx),
				client.OriginalMode(true),
				client.DebugMode(runClientDebug),
				client.ModelName(runClientModels),
				client.IterationCount(1),
				client.ProfileMemory(true),
				client.UploadProfile(false),
				client.ConcurrentRunCount(1),
				client.ShowProgress(false),
				client.ProfileIO(false),
				client.PostProcess(true),
				client.EagerInitialize(true),
				client.Stdout(ioutil.Discard),
				client.Stderr(ioutil.Discard),
				client.LargeModels(runClientLargeModels),
			}
		if runClientDebug {
			opts = append(opts, []client.Option{
				client.Stdout(os.Stdout),
				client.Stderr(os.Stderr),
			}...)
		}
		client := client.New(opts...)
		traces, err := client.Run(true)
		if err != nil {
			return err
		}
		//pp.Println(traces)
		meminfo := trace.MemoryInformation(traces)
		meminfo.Write(runClientMemoryOutputFormat, os.Stdout)
		return err
	},
}

func init() {
	clientRunMemoryCmd.Flags().StringVar(&runClientModels, "models", "all", "List of models to use (comma seperated)")
	clientRunMemoryCmd.Flags().BoolVarP(&runClientDebug, "debug", "d", false, "Print debug messages from the client")
	clientRunMemoryCmd.Flags().StringVarP(&runClientMemoryOutputFormat, "format", "f", "table", "Output format to print the memory information")
	clientRunMemoryCmd.Flags().BoolVar(&runClientLargeModels, "large_models", false, "run the large models")
	clientCmd.AddCommand(clientRunMemoryCmd)
}

package cmd

import (
	"context"
	"os"

	"github.com/rai-project/micro18-tools/pkg/client"
	"github.com/spf13/cobra"
)

var clientRunMemoryCmd = &cobra.Command{
	Use:     "memory_usage",
	Aliases: []string{"memory"},
	Short:   "Run the client command and print out how much memory is used by internal layers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client := client.New(
			client.Context(ctx),
			client.ModelName(runClientModels),
			client.IterationCount(1),
			client.ProfileMemory(true),
			client.UploadProfile(false),
		)
		traces, err := client.Run()
		if err != nil {
			return err
		}
		meminfo := trace.MemoryInformation{traces}
		meminfo.Write("table", os.Stdout)
		return err
	},
}

func init() {
	clientRunMemoryCmd.Flags().StringVar(&runClientModels, "models", "all", "List of models to use (comma seperated)")
	clientCmd.AddCommand(clientRunMemoryCmd)
}

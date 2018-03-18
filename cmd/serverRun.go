package cmd

import (
	"context"

	"github.com/rai-project/micro18-tools/pkg/server"
	"github.com/spf13/cobra"
)

var (
	runServerDebug               bool
	runServerEvictionPolicy      string
	runServerModelEstimationRate float32
	runServerMemoryPercentage    float32
	runServerPersistCPU          bool
)

func makeServerRun(ctx context.Context) *server.Server {
	return server.New(
		server.Context(ctx),
		server.DebugMode(runServerDebug),
		server.EvictionPolicy(runServerEvictionPolicy),
		server.ModelEstimationRate(runServerModelEstimationRate),
		server.MemoryPercentage(runServerMemoryPercentage),
		server.PersistCPU(runServerPersistCPU),
	)
}

// serverRunCmd represents the serverRun command
var serverRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the server command and produce profile files",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if ok, err := server.IsValidEvictionPolicy(runServerEvictionPolicy); !ok {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		server := makeServerRun(ctx)
		_, err := server.Run()
		return err
	},
}

func init() {
	serverRunCmd.Flags().BoolVarP(&runServerDebug, "debug", "d", false, "Print debug messages from the client")
	serverRunCmd.Flags().StringVar(&runServerEvictionPolicy, "eviction", "lru", "Eviction policy used by the server")
	serverRunCmd.Flags().Float32Var(&runServerModelEstimationRate, "model_estimation_rate", 1.0, "File size multiplier used to determine how much memory would be used by a model")
	serverRunCmd.Flags().Float32Var(&runServerMemoryPercentage, "memory_percentage", 0.8, "Percentage of GPU memory that can be used to persist models")
	serverRunCmd.Flags().BoolVar(&runServerPersistCPU, "persist_cpu", true, "Persist memory on CPU to avoid rereading the data from disk after eviction")

	serverCmd.AddCommand(serverRunCmd)
}

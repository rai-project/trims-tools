package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/Unknwon/com"

	"github.com/rai-project/uuid"

	mconfig "github.com/rai-project/micro18-tools/pkg/config"
	"github.com/rai-project/micro18-tools/pkg/server"
	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/spf13/cobra"
)

var (
	runServerID                  string
	runServerDebug               bool
	runServerEvictionPolicy      string
	runServerModelEstimationRate float32
	runServerMemoryPercentage    float32
	runServerPersistCPU          bool
	runServerWriteProfile        bool
	serverInfoPath               string
)

func makeServerRun(ctx context.Context) *server.Server {
	return server.New(
		server.Context(ctx),
		server.ID(runServerID),
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

		runServerID = uuid.NewV4()
		serverInfoPath = mconfig.Config.ServerInfoPath

		if ok, err := server.IsValidEvictionPolicy(runServerEvictionPolicy); !ok {
			return err
		}
		if server.IsRunning() {
			return errors.New("the server is already running. make sure to kill it")
		}

		info := trace.TraceServerInfo{
			ID:               runServerID,
			StartTime:        time.Now(),
			EvictionPolicty:  runServerEvictionPolicy,
			EstimationRate:   runServerModelEstimationRate,
			MemoryPercentage: runServerMemoryPercentage,
			PersistCPU:       runServerPersistCPU,
		}
		if bts, err := json.Marshal(info); err == nil {
			ioutil.WriteFile(serverInfoPath, bts, 0644)
		}
		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		if com.IsFile(serverInfoPath) {
			os.Remove(serverInfoPath)
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
	serverRunCmd.Flags().BoolVar(&runServerWriteProfile, "profile", false, "Write the server profile file")
	serverRunCmd.Flags().StringVar(&runServerEvictionPolicy, "eviction", "lru", "Eviction policy used by the server")
	serverRunCmd.Flags().Float32Var(&runServerModelEstimationRate, "model_estimation_rate", 1.0, "File size multiplier used to determine how much memory would be used by a model")
	serverRunCmd.Flags().Float32Var(&runServerMemoryPercentage, "memory_percentage", 0.8, "Percentage of GPU memory that can be used to persist models")
	serverRunCmd.Flags().BoolVar(&runServerPersistCPU, "persist_cpu", true, "Persist memory on CPU to avoid rereading the data from disk after eviction")

	serverCmd.AddCommand(serverRunCmd)
}

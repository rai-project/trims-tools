package cmd

import (
	"context"

	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/spf13/cobra"
)

var (
	runClientNTimes int
)

// runCmd represents the run command
var clientRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the client command and produce profile files",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		for ii := 0; ii < runClientNTimes; ii++ {
			trace.Run(ctx, ii)
		}
	},
}

func init() {
	clientCmd.AddCommand(clientRunCmd)
	clientRunCmd.Flags().IntVarP(&runClientNTimes, "iterations", "n", 1, "Number of iterations to run the client")
}

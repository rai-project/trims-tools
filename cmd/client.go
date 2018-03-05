package cmd

import (
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use: "client",
}

func init() {
	rootCmd.AddCommand(clientCmd)
}

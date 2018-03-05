package cmd

import (
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use: "download",
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}

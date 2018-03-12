package cmd

import (
	"context"

	"github.com/rai-project/micro18-tools/pkg/assets"
	"github.com/spf13/cobra"
)

// downloadModelCmd represents the downloadModel command
var downloadAssetsCmd = &cobra.Command{
	Use:     "assets",
	Aliases: []string{"model", "models", "feature", "features"},
	Short:   "Downloads assets to directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		return assets.Download(context.Background())
	},
}

func init() {
	downloadCmd.AddCommand(downloadAssetsCmd)
}

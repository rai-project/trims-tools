package cmd

import (
	"fmt"

	"github.com/rai-project/micro18-tools/pkg/config"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var infoConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "shows the config paramters in yml format",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.Config
		bts, err := yaml.Marshal(c)
		if err != nil {
			return err
		}
		fmt.Println(string(bts))
		return nil
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "shows some information about the tool",
}

func init() {
	infoCmd.AddCommand(infoConfigCmd)
	rootCmd.AddCommand(infoCmd)
}

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/spf13/cobra"
)

var traceMemoryInfoCmd = &cobra.Command{
	Use:     "memoryinfo",
	Aliases: []string{"memory"},
	Short:   "Shows amount of memory allocated via cudaMalloc in trace",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, path := range args {
			if !com.IsFile(path) {
				return errors.Errorf("the profile file %s was not found", path)
			}
			bts, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrapf(err, "unable to read the profile file from %s", path)
			}
			var tr trace.Trace
			if err := json.Unmarshal(bts, &tr); err != nil {
				return errors.Wrapf(err, "unable to unmarshal the profile file from %s", path)
			}
			meminfo := trace.MemoryInformation([]*trace.Trace{&tr})
			meminfo.Write(runClientMemoryOutputFormat, os.Stdout)
		}
		return nil
	},
}

func init() {
	traceCmd.AddCommand(traceMemoryInfoCmd)
}

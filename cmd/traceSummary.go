package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/spf13/cobra"
)

var (
	traceSummarizeOutputFile string
)

// traceSummarizeCmd represents the traceSummarize command
var traceSummarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "A brief description of your command",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.Errorf("not profiles were specified")
		}
		var res []*trace.TraceSummary
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
			ts, err := tr.Summarize()
			if err != nil {
				return err
			}
			res = append(res, ts)
		}
		bts, err := json.Marshal(res)
		if err != nil {
			return errors.Wrap(err, "unable to marshal query results")
		}
		if !com.IsDir(filepath.Dir(traceSummarizeOutputFile)) {
			os.MkdirAll(filepath.Dir(traceSummarizeOutputFile), os.ModePerm)
		}
		if err := ioutil.WriteFile(traceSummarizeOutputFile, bts, 0644); err != nil {
			return errors.Wrapf(err, "unable to write query results to %s", traceSummarizeOutputFile)
		}
		fmt.Println("output is written to " + traceSummarizeOutputFile)
		return nil
	},
}

func init() {
	traceCmd.AddCommand(traceSummarizeCmd)
	traceSummarizeCmd.Flags().StringVarP(&traceSummarizeOutputFile, "output", "o", "summary.json", "Ther output path to the trace summary")
}

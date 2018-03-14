// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	traceCombineOutputFile string
	traceCombineAdjust     bool
	traceCombineSkipFirst  bool
)

// traceCombineCmd represents the traceCombine command
var traceCombineCmd = &cobra.Command{
	Use:   "combine",
	Short: "A brief description of your command",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		traces := []trace.Trace{}
		if traceCombineSkipFirst {
			args = args[1:]
		}
		if len(args) == 0 {
			return errors.Errorf("not profiles were specified. check to see if you skipped the first profile")
		}
		for _, path := range args {
			if !com.IsFile(path) {
				return errors.Errorf("the profile file %s was not found", path)
			}
			bts, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrapf(err, "unable to read the profile file from %s", path)
			}
			var trace trace.Trace
			if err := json.Unmarshal(bts, &trace); err != nil {
				return errors.Wrapf(err, "unable to unmarshal the profile file from %s", path)
			}
			if traceCombineAdjust {
				trace, err = trace.Adjust()
				if err != nil {
					log.WithError(err).Error("failed to adjust trace")
				}
			}
			traces = append(traces, trace)
		}
		combinedTrace := trace.Combine(traces[0], traces[1:]...)
		bts, err := json.Marshal(combinedTrace)
		if err != nil {
			return errors.Wrap(err, "unable to marshal combined traces")
		}
		if !com.IsDir(filepath.Dir(traceCombineOutputFile)) {
			os.MkdirAll(filepath.Dir(traceCombineOutputFile), os.ModePerm)
		}
		if err := ioutil.WriteFile(traceCombineOutputFile, bts, 0644); err != nil {
			return errors.Wrapf(err, "unable to write combined traces to %s", traceCombineOutputFile)
		}
		fmt.Println("output is written to " + traceCombineOutputFile)
		return nil
	},
}

func init() {
	traceCmd.AddCommand(traceCombineCmd)
	traceCombineCmd.Flags().StringVarP(&traceCombineOutputFile, "output", "o", "combined.json", "Ther output path to the combined trace")
	traceCombineCmd.Flags().BoolVar(&traceCombineAdjust, "adjust", true, "Adjust the timeline to ignore categories, adjust event names, and zero out the trace")
	traceCombineCmd.Flags().BoolVar(&traceCombineSkipFirst, "skip_first", false, "Skip the first input argument")
}

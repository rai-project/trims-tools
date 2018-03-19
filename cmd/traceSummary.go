package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/Jeffail/tunny"
	"github.com/Unknwon/com"
	zglob "github.com/mattn/go-zglob"
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
	Short: "Summarizes the traces within a directory or list of files",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var res []*trace.TraceSummary
		files := []string{}
		for _, path := range args {
			if com.IsDir(path) {
				matches, err := zglob.Glob(filepath.Join(path, "**", "*.json"))
				if err == nil {
					files = append(files, matches...)
				}
				// matches, err = zglob.Glob(filepath.Join(path, "*.json"))
				// if err == nil {
				// 	files = append(files, matches...)
				// }
			} else {
				files = append(files, path)
			}
		}
		var mut sync.Mutex
		var wg sync.WaitGroup

		processFile := func(path0 interface{}) interface{} {
			defer wg.Done()
			path := path0.(string)
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
			mut.Lock()
			defer mut.Unlock()
			res = append(res, ts)
			return err
		}

		processPool := tunny.NewFunc(runtime.NumCPU(), processFile)
		defer processPool.Close()

		for _, path := range files {
			if !com.IsFile(path) {
				return errors.Errorf("the profile file %s was not found", path)
			}
			wg.Add(1)
			go func(path string) {
				processPool.Process(path)
			}(path)
		}
		wg.Wait()
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
	traceSummarizeCmd.Flags().StringVarP(&traceSummarizeOutputFile, "output", "o", "summary.json", "The output path to the trace summary")
}

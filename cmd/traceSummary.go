package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/Jeffail/tunny"
	"github.com/Unknwon/com"
	zglob "github.com/mattn/go-zglob"
	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/rai-project/micro18-tools/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	traceSummarizeOutputFile string
	traceSummarizeDetailed   bool
	traceSummarizeDeep       bool
)

func traceSummarize(cmd *cobra.Command, args []string) error {
	outputFile := traceSummarizeOutputFile
	if outputFile == "" {
		fst := args[0]
		if com.IsDir(fst) {
			outputFile = filepath.Join(fst, "summary.json")
		} else {
			if cwd, err := os.Getwd(); err != nil {
				outputFile = filepath.Join(cwd, "summary.json")
			} else {
				outputFile = "summary.json"
			}
		}
	}

	var res []*trace.TraceSummary
	files := []string{}
	for _, path := range args {
		if com.IsDir(path) {
			matches, err := zglob.Glob(filepath.Join(path, "**", "*.json"))
			if err == nil {
				for _, m := range matches {
					if m == outputFile {
						continue
					}
					files = append(files, m)
				}
			}
			// matches, err = zglob.Glob(filepath.Join(path, "*.json"))
			// if err == nil {
			// 	files = append(files, matches...)
			// }
		} else if path != outputFile {
			files = append(files, path)
		}
	}

	var mut sync.Mutex
	var wg sync.WaitGroup

	progress := utils.NewProgress("summarizing traces", len(files))
	defer progress.FinishPrint("finished summarizing traces and places the result in " + outputFile)

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
		ts, err := tr.Summarize(traceSummarizeDetailed)
		if err != nil {
			return err
		}
		mut.Lock()
		defer mut.Unlock()
		progress.Increment()
		res = append(res, ts)
		return err
	}

	processPool := tunny.NewFunc(2*runtime.NumCPU(), processFile)
	defer processPool.Close()

	for _, path := range files {
		if !com.IsFile(path) {
			return errors.Errorf("the profile file %s was not found", path)
		}
		if strings.Contains(path, "combined-") {
			continue
		}
		if strings.Contains(path, "compared-") {
			continue
		}
		if strings.Contains(path, "summary-") {
			continue
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
	if !com.IsDir(filepath.Dir(outputFile)) {
		os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
	}
	if err := ioutil.WriteFile(outputFile, bts, 0644); err != nil {
		return errors.Wrapf(err, "unable to write query results to %s", outputFile)
	}
	return nil
}

// traceSummarizeCmd represents the traceSummarize command
var traceSummarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Summarizes the traces within a directory or list of files",
	Long: "example usage\n" +
		"go run main.go trace summarize ~/micro18_profiles/`hostname`/Exponential* ~/micro18_profiles/`hostname`/Pareto* ~/micro18_profiles/`hostname`/Weibull* ~/micro18_profiles/`hostname`/Poisson* ~/micro18_profiles/`hostname`/Uniform*\n" +
		"or\n" +
		"go run main.go trace summarize ~/micro18_profiles/`hostname`/*Exponential_rt* ~/micro18_profiles/`hostname`/*Pareto_xm1_l* ~/micro18_profiles/`hostname`/*Weibull_k* ~/micro18_profiles/`hostname`/*Poisson_l* ~/micro18_profiles/`hostname`/*Uniform_min0_max1*\n" +
		"",
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if traceSummarizeDeep {
			return traceSummarize(cmd, args)
		}
		for _, a := range args {
			if !com.IsDir(a) {
				return errors.New("the directory " + a + " was not found. expected a set of directories if the option --deep=true is passed in")
			}
		}

		for _, a := range args {
			fmt.Println("summarizing ", a)
			err := traceSummarize(cmd, []string{a})
			if err != nil {
				log.WithError(err).WithField("directory", a).Error("failed to summarize trace")
			}
		}
		return nil
	},
}

func init() {
	traceCmd.AddCommand(traceSummarizeCmd)
	traceSummarizeCmd.Flags().StringVarP(&traceSummarizeOutputFile, "output", "o", "", "The output path to the trace summary")
	traceSummarizeCmd.Flags().BoolVar(&traceSummarizeDetailed, "detailed", true, "The output should contain more detailed event information")
	traceSummarizeCmd.Flags().BoolVar(&traceSummarizeDeep, "deep", false, "If enabled and the input is a set of directories, then the summary.json is for all the json files for all directory. Otherwise a summary.json is generated for each directory")
}

package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Jeffail/tunny"
	"github.com/rai-project/micro18-tools/pkg/workload"

	"github.com/rai-project/uuid"

	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/assets"
	"github.com/rai-project/micro18-tools/pkg/config"
	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/rai-project/micro18-tools/pkg/utils"
)

func (c Client) Run() ([]*trace.Trace, error) {
	options := c.options

	if options.modelDistribution == "none" {
		return c.run()
	}
	return c.runWorkload()
}

func (c Client) run() ([]*trace.Trace, error) {
	options := c.options

	models := assets.Models
	if strings.ToLower(options.modelName) != "all" {
		models = assets.ModelManifests{}
		for _, m := range assets.Models {
			if strings.ToLower(m.MustCanonicalName()) == strings.ToLower(options.modelName) {
				models = assets.ModelManifests{m}
				break
			}
		}
		if len(models) == 0 {
			return nil, errors.Errorf("the model %s was not found in the asset list", options.modelName)
		}
	}

	progress := utils.NewProgress("running client models", len(models)*options.iterationCount)
	defer progress.FinishPrint("finished running client")

	var res []*trace.Trace
	var combined *trace.Trace
	var mut sync.Mutex
	var wg sync.WaitGroup
	iterationCount := map[string]int{}

	runModel := func(arg interface{}) interface{} {
		model := arg.(assets.ModelManifest)
		defer progress.Increment()
		trace, err := c.RunOnce(model)
		if err != nil {
			return nil
		}
		if trace == nil {
			return nil
		}
		go func() {
			defer wg.Done()
			mut.Lock()
			defer mut.Unlock()
			cannonicalName := model.MustCanonicalName()
			ii, ok := iterationCount[cannonicalName]
			if !ok {
				ii = 0
				iterationCount[cannonicalName] = ii
			} else {
				iterationCount[cannonicalName] = ii + 1
			}
			trace.Iteration = int64(ii)
			if combined == nil {
				combined = trace
				combined.ID = uuid.NewV4()
			} else {
				combined.Combine(*trace)
			}
			if combined != nil {
				if err := combined.Upload(); err != nil {
					log.WithError(err).Error("failed to upload combined profile output")
				}
				res = append(res, combined)
			}
		}()
		return nil
	}

	execPool := tunny.NewFunc(options.concurrentRunCount, runModel)

	for _, model := range models {
		// progress.Prefix(fmt.Sprintf("running client model %s", model.MustCanonicalName()))
		for ii := 0; ii < options.iterationCount; ii++ {
			wg.Add(1)
			//progress.Prefix(fmt.Sprintf("running client model %s", model.MustCanonicalName()))
			go func(model assets.ModelManifest) {
				execPool.Process(model)
			}(model)
		}
	}
	wg.Wait()

	if combined != nil {
		if err := combined.Upload(); err != nil {
			log.WithError(err).Error("failed to upload combined profile output")
		}
		res = append(res, combined)
	}
	return res, nil
}

func (c Client) runWorkload() ([]*trace.Trace, error) {
	options := c.options

	models := assets.Models
	if strings.ToLower(options.modelName) != "all" {
		models = assets.ModelManifests{}
		for _, m := range assets.Models {
			if strings.ToLower(m.MustCanonicalName()) == strings.ToLower(options.modelName) {
				models = assets.ModelManifests{m}
				break
			}
		}
		if len(models) == 0 {
			return nil, errors.Errorf("the model %s was not found in the asset list", options.modelName)
		}
	}
	modelGen, err := workload.New(options.modelDistribution, options.modelDistributionParams)
	if err != nil {
		return nil, err
	}

	progress := utils.NewProgress("running client models", len(models)*options.iterationCount)
	defer progress.FinishPrint("finished running client")

	var res []*trace.Trace
	var combined *trace.Trace
	var mut sync.Mutex
	var wg sync.WaitGroup
	iterationCount := map[string]int{}

	runModel := func(arg interface{}) interface{} {
		defer progress.Increment()

		model := arg.(assets.ModelManifest)
		cannonicalName := model.MustCanonicalName()

		mut.Lock()
		iterCnt, ok := iterationCount[cannonicalName]
		if !ok {
			iterCnt = 1
			iterationCount[cannonicalName] = iterCnt
		} else {
			iterationCount[cannonicalName] = iterCnt + 1
		}
		mut.Unlock()

		if options.modelIterationCount != -1 && iterCnt >= options.modelIterationCount {
			return nil
		}
		trace, err := c.RunOnce(model)
		if err != nil {
			return nil
		}
		if trace == nil {
			return nil
		}

		trace.Iteration = int64(iterCnt)

		go func() {
			defer wg.Done()
			mut.Lock()
			defer mut.Unlock()
			if combined == nil {
				combined = trace
				combined.ID = uuid.NewV4()
			} else {
				combined.Combine(*trace)
			}
			if combined != nil {
				if err := combined.Upload(); err != nil {
					log.WithError(err).Error("failed to upload combined profile output")
				}
				res = append(res, combined)
			}
		}()
		return nil
	}

	execPool := tunny.NewFunc(options.concurrentRunCount, runModel)

	for model := range modelGen.ModelGenerator(models) {
		wg.Add(1)
		//progress.Prefix(fmt.Sprintf("running client model %s", model.MustCanonicalName()))
		go func(model assets.ModelManifest) {
			execPool.Process(model)
		}(model)
	}
	wg.Wait()

	if combined != nil {
		if err := combined.Upload(); err != nil {
			log.WithError(err).Error("failed to upload combined profile output")
		}
		res = append(res, combined)
	}
	return res, nil
}

func (c Client) RunOnce(model assets.ModelManifest) (*trace.Trace, error) {
	options := c.options

	dims, err := model.GetImageDimensions()
	if err != nil {
		dims = []uint32{3, 224, 224}
	}
	mean, err := model.GetMeanImage()
	if err != nil {
		mean = []float32{0, 0, 0}
	}
	if len(dims) != 3 {
		err := errors.Errorf("expecting a 3 element vector for dimensions %v", dims)
		return nil, err
	}
	cannonicalName := model.MustCanonicalName()

	id := uuid.NewV4()
	profileFilePath := filepath.Join(config.Config.ProfileOutputDirectory, time.Now().Format("Jan-_2-15"), fmt.Sprintf("%s_%s.json", cannonicalName, id))
	env := map[string]string{
		"UPR_ENABLED":                 "true",
		"UPR_RUN_ID":                  id,
		"DATE":                        time.Now().Format(time.RFC3339Nano),
		"UPR_MODEL_NAME":              cannonicalName,
		"UPR_CLIENT":                  "1",
		"MXNET_CPU_PRIORITY_NTHREADS": "1",
		"OMP_NUM_THREADS":             "1",
		"MXNET_ENGINE_TYPE":           "NaiveEngine",
		"MXNET_GPU_WORKER_NTHREADS":   "1",
		"UPR_PROFILE_TARGET":          profileFilePath,
		"UPR_INPUT_CHANNELS":          strconv.Itoa(int(dims[0])),
		"UPR_INPUT_HEIGHT":            strconv.Itoa(int(dims[1])),
		"UPR_INPUT_WIDTH":             strconv.Itoa(int(dims[2])),
		"UPR_INPUT_MEAN_R":            fmt.Sprintf("%v", mean[0]),
		"UPR_INPUT_MEAN_G":            fmt.Sprintf("%v", mean[1]),
		"UPR_INPUT_MEAN_B":            fmt.Sprintf("%v", mean[2]),
		"UPR_GIT_SHA":                 config.Version.GitCommit,
		"UPR_GIT_BRANCH":              config.Version.GitBranch,
		"UPR_GIT_Date":                config.Version.BuildDate,
	}
	if options.original {
		env["UPR_ENABLED"] = "true"
	} else {
		env["UPR_ENABLED"] = "false"
	}
	if options.debug {
		env["GLOG_logtostderr"] = "1"
		env["GLOG_v"] = "0"
		env["GLOG_stderrthreshold"] = "0"
	}
	if options.eagerInitialize {
		env["UPR_INITIALIZE_EAGER"] = "true"
	}
	if options.eagerInitializeAsync {
		env["UPR_INITIALIZE_EAGER_ASYNC"] = "true"
	}
	ran, err := utils.ExecCmd(
		config.Config.ClientPath,
		env,
		os.Stdout,
		os.Stderr,
		config.Config.ClientRunCmd,
	)
	if !ran {
		err := errors.Errorf("failed to run cmd %s", config.Config.ClientRunCmd)
		log.WithError(err).Error("failed to run model")
		return nil, err
	}
	if err != nil {
		log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to run model")
		return nil, err
	}
	if !options.postprocess {
		return nil, nil
	}
	bts, err := ioutil.ReadFile(profileFilePath)
	if err != nil {
		err = errors.Wrapf(err, "unable to read profile file %s", profileFilePath)
		log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to read profile output")
		return nil, err
	}
	var trace trace.Trace
	if err := json.Unmarshal(bts, &trace); err != nil {
		err = errors.Wrapf(err, "unable to unmarshal profile file %s", profileFilePath)
		log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to unmarshal profile output")
		return nil, err
	}
	if options.uploadProfile {
		if err := trace.Upload(); err != nil {
			err = errors.Wrapf(err, "unable to upload profile file %s", profileFilePath)
			log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to upload profile output")
		}
	}
	return &trace, nil
}

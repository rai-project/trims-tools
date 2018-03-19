package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"github.com/Unknwon/com"
	"github.com/fatih/color"
	"github.com/k0kubun/pp"
	"github.com/spf13/cast"

	"github.com/Jeffail/tunny"
	"github.com/cheggaaa/pb"
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

	cmdPath := filepath.Join(config.Config.ClientPath, config.Config.ClientRunCmd)
	if !com.IsFile(cmdPath) {
		return nil, errors.Errorf("the client command %s was not found in %s. make sure that the code compiled correctly",
			config.Config.ClientRunCmd, config.Config.ClientPath)
	}

	fmt.Println(color.GreenString("✱ Running client and placing profile in " + config.Config.ProfileOutputDirectory))
	if com.IsFile(config.Config.ServerInfoPath) && config.Config.UPREnabled {
		bts, err := ioutil.ReadFile(config.Config.ServerInfoPath)
		if err == nil {
			var info trace.TraceServerInfo
			if err := json.Unmarshal(bts, &info); err == nil {
				fmt.Println(color.GreenString("✱ The server id for the run is " + info.ID))
			}
		}
	}

	if options.modelDistribution == "none" || options.modelDistribution == "" {
		return c.run()
	}
	return c.runWorkload()
}

func (c Client) run() ([]*trace.Trace, error) {
	options := c.options

	models, err := assets.FilterModels(options.modelName)
	if err != nil {
		return nil, err
	}

	if options.showProgress && len(models) <= 1 {
		options.showProgress = false
	}

	var progress *pb.ProgressBar
	if options.showProgress {
		progress = utils.NewProgress("running client models", len(models)*options.iterationCount)
		defer progress.FinishPrint("finished running client")
	}

	var res []*trace.Trace
	var combined *trace.Trace
	var mut sync.Mutex
	var wg sync.WaitGroup
	iterationCount := map[string]int{}

	runModel := func(arg interface{}) interface{} {
		defer wg.Done()
		model := arg.(assets.ModelManifest)
		if options.showProgress {
			defer progress.Increment()
		}
		profilePath, duration, err := c.RunOnce(model)
		if err != nil {
			return nil
		}
		if profilePath == "" {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			mut.Lock()
			defer mut.Unlock()

			trace, err := c.readProfile(profilePath, duration)
			if err != nil {
				return
			}
			if trace == nil {
				return
			}

			cannonicalName := model.MustCanonicalName()
			iterCnt, ok := iterationCount[cannonicalName]
			if !ok {
				iterCnt = 0
				iterationCount[cannonicalName] = iterCnt
			} else {
				iterationCount[cannonicalName] = iterCnt + 1
			}
			trace.Iteration = int64(iterCnt)
			trace.OtherDataRaw.Iteration = int64(iterCnt)
			if trace.OtherDataRaw != nil {
				pp.Println(trace.OtherDataRaw.EndToEndProcessTime)
			}
			trace.OtherData = nil
			if bts, err := json.Marshal(trace); err == nil {
				ioutil.WriteFile(profilePath, bts, 0644)
			}

			res = append(res, trace)

			if combined == nil {
				combined = trace
				combined.ID = uuid.NewV4()
			} else {
				combined.Combine(*trace)
			}
		}()
		return nil
	}

	execPool := tunny.NewFunc(options.concurrentRunCount, runModel)
	defer execPool.Close()

	for _, model := range models {
		// if options.showProgress {
		// progress.Prefix(fmt.Sprintf("running client model %s", model.MustCanonicalName()))
		// }
		for ii := 0; ii < options.iterationCount; ii++ {
			wg.Add(1)
			// if options.showProgress {
			//progress.Prefix(fmt.Sprintf("running client model %s", model.MustCanonicalName()))
			// }
			go func(model assets.ModelManifest) {
				execPool.Process(model)
			}(model)
		}
	}
	wg.Wait()

	if combined != nil {
		// online one combined result
		if len(combined.OtherData) <= 1 {
			combined = nil
		}
	}
	if combined != nil {
		id := uuid.NewV4()
		path := filepath.Join(config.Config.ProfileOutputDirectory, "combined-"+id+".json")
		bts, err := json.Marshal(combined)
		if err == nil {
			ioutil.WriteFile(path, bts, 0644)
		}
	}
	if combined != nil && options.uploadProfile {
		if err := combined.Upload(); err != nil {
			log.WithError(err).Error("failed to upload combined profile output")
		}
	}
	if combined != nil {
		res = append(res, combined)
	}
	return res, nil
}

func (c Client) runWorkload() ([]*trace.Trace, error) {
	options := c.options

	models, err := assets.FilterModels(options.modelName)
	if err != nil {
		return nil, err
	}

	if options.showProgress && len(models) <= 1 {
		options.showProgress = false
	}

	modelGen, err := workload.New(options.modelDistribution, options.modelDistributionParams)
	if err != nil {
		return nil, err
	}

	var progress *pb.ProgressBar
	if options.showProgress {
		progress = utils.NewProgress("running client models", options.iterationCount)
		defer progress.FinishPrint("finished running client")
	}

	var res []*trace.Trace
	var combined *trace.Trace
	var mut sync.Mutex
	var wg sync.WaitGroup
	iterationCount := map[string]int{}

	runModel := func(arg interface{}) interface{} {
		defer wg.Done()
		if options.showProgress {
			defer progress.Increment()
		}

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
		profilePath, duration, err := c.RunOnce(model)
		if err != nil {
			return nil
		}
		if profilePath == "" {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			mut.Lock()
			defer mut.Unlock()

			trace, err := c.readProfile(profilePath, duration)
			if err != nil {
				return
			}
			if trace == nil {
				return
			}

			trace.Iteration = int64(iterCnt)
			trace.OtherDataRaw.Iteration = int64(iterCnt)
			if trace.OtherDataRaw != nil {
				pp.Println(trace.OtherDataRaw.EndToEndProcessTime)
			}
			trace.OtherData = nil
			if bts, err := json.Marshal(trace); err == nil {
				ioutil.WriteFile(profilePath, bts, 0644)
			}
			res = append(res, trace)

			if combined == nil {
				combined = trace
				combined.ID = uuid.NewV4()
			} else {
				combined.Combine(*trace)
			}
		}()
		return nil
	}

	execPool := tunny.NewFunc(options.concurrentRunCount, runModel)
	defer execPool.Close()

	ii := 0

	for model := range modelGen.ModelGenerator(models) {
		if ii == options.iterationCount {
			break
		}
		ii++
		wg.Add(1)
		// if options.showProgress {
		//progress.Prefix(fmt.Sprintf("running client model %s", model.MustCanonicalName()))
		// }
		go func(model assets.ModelManifest) {
			defer func() {
				if r := recover(); r != nil {
					return
				}
			}()
			execPool.Process(model)
		}(model)
	}
	wg.Wait()

	if combined != nil {
		// online one combined result
		if len(combined.OtherData) <= 1 {
			combined = nil
		}
	}
	if combined != nil {
		id := uuid.NewV4()
		path := filepath.Join(config.Config.ProfileOutputDirectory, "combined-"+id+".json")
		bts, err := json.Marshal(combined)
		if err == nil {
			ioutil.WriteFile(path, bts, 0644)
		}
	}
	if combined != nil && options.uploadProfile {
		if err := combined.Upload(); err != nil {
			log.WithError(err).Error("failed to upload combined profile output")
		}
	}
	if combined != nil {
		res = append(res, combined)
	}
	return res, nil
}

func (c Client) RunOnce(model assets.ModelManifest) (string, time.Duration, error) {
	options := c.options

	dims, err := model.GetImageDimensions()
	if err != nil {
		log.WithError(err).Errorf("failed to get image dimensions for %s", model.MustCanonicalName())
		dims = []uint32{3, 224, 224}
	}
	mean, err := model.GetMeanImage()
	if err != nil {
		mean = []float32{0, 0, 0}
	}
	if len(dims) != 3 {
		err := errors.Errorf("expecting a 3 element vector for dimensions %v", dims)
		return "", 0, err
	}
	cannonicalName := model.MustCanonicalName()

	id := uuid.NewV4()
	profileBaseName := fmt.Sprintf("%s_%s.json", cannonicalName, id)
	if !options.original {
		profileBaseName = "upr_" + profileBaseName
	}
	profileFilePath := filepath.Join(config.Config.ProfileOutputDirectory, profileBaseName)
	env := map[string]string{
		"HOME":                         config.HomeDir,
		"UPR_RUN_ID":                   id,
		"DATE":                         time.Now().Format(time.RFC3339Nano),
		"UPR_MODEL_NAME":               cannonicalName,
		"UPR_CLIENT":                   "1",
		"MXNET_CPU_PRIORITY_NTHREADS":  "1",
		"OMP_NUM_THREADS":              "1",
		"MXNET_ENGINE_TYPE":            "NaiveEngine",
		"MXNET_GPU_WORKER_NTHREADS":    "1",
		"UPR_BASE_DIR":                 config.Config.BasePath + "/",
		"UPR_PROFILE_TARGET":           profileFilePath,
		"UPR_INPUT_CHANNELS":           cast.ToString(dims[0]),
		"UPR_INPUT_HEIGHT":             cast.ToString(dims[1]),
		"UPR_INPUT_WIDTH":              cast.ToString(dims[2]),
		"UPR_INPUT_MEAN_R":             fmt.Sprintf("%v", mean[0]),
		"UPR_INPUT_MEAN_G":             fmt.Sprintf("%v", mean[1]),
		"UPR_INPUT_MEAN_B":             fmt.Sprintf("%v", mean[2]),
		"UPR_GIT_SHA":                  config.Version.GitCommit,
		"UPR_GIT_BRANCH":               config.Version.GitBranch,
		"UPR_GIT_Date":                 config.Version.BuildDate,
		"CUDA_VISIBLE_DEVICES":         config.Config.VisibleDevices,
		"MXNET_CUDNN_AUTOTUNE_DEFAULT": "0",
	}
	if options.original {
		env["UPR_ENABLED"] = "false"
	} else {
		env["UPR_ENABLED"] = "true"
	}
	if options.profileMemory {
		env["UPR_ENABLE_MEMORY_PROFILE"] = "true"
	} else {
		env["UPR_ENABLE_MEMORY_PROFILE"] = "false"
	}
	if options.debug {
		env["GLOG_logtostderr"] = "1"
		env["GLOG_v"] = "0"
		env["GLOG_stderrthreshold"] = "0"
	}
	if options.eagerInitialize {
		env["UPR_INITIALIZE_EAGER"] = "true"
	} else {
		env["UPR_INITIALIZE_EAGER"] = "false"
	}
	if options.eagerInitializeAsync {
		env["UPR_INITIALIZE_EAGER_ASYNC"] = "true"
	} else {
		env["UPR_INITIALIZE_EAGER_ASYNC"] = "false"
	}
	if false {
		pp.Println(env)
	}
	tic := time.Now()
	ran, err := utils.ExecCmd(
		config.Config.ClientPath,
		env,
		options.stdout,
		options.stderr,
		config.Config.ClientRunCmd,
	)
	timeToRun := time.Since(tic)
	if !ran {
		path := filepath.Join(config.Config.ClientPath, config.Config.ClientRunCmd)
		err := errors.Errorf("failed to run cmd %s", path)
		log.WithError(err).WithField("model_name", cannonicalName).WithField("dims", dims).Error("failed to run model")
		return "", 0, err
	}
	if err != nil {
		path := filepath.Join(config.Config.ClientPath, config.Config.ClientRunCmd)
		err = errors.Wrapf(err, "failed to run cmd %s", path)
		lg := log.WithField("cmd", config.Config.ClientRunCmd).WithField("model_name", cannonicalName).WithField("dims", dims)
		if options.debug {
			lg = lg.WithError(err)
		}
		lg.Error("failed to run model")
		return "", 0, err
	}
	if !options.postprocess {
		return "", timeToRun, nil
	}
	return profileFilePath, timeToRun, nil
}

func (c *Client) readProfile(profileFilePath string, timeToRun time.Duration) (*trace.Trace, error) {
	bts, err := ioutil.ReadFile(profileFilePath)
	if err != nil {
		err = errors.Wrapf(err, "unable to read profile file %s", profileFilePath)
		log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to read profile output")
		return nil, err
	}
	trace := new(trace.Trace)
	if err := json.Unmarshal(bts, trace); err != nil {
		err = errors.Wrapf(err, "unable to unmarshal profile file %s", profileFilePath)
		log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to unmarshal profile output")
		return nil, err
	}
	if trace.OtherDataRaw != nil {
		trace.OtherDataRaw.EndToEndProcessTime = timeToRun
	}
	if len(trace.OtherData) != 0 {
		for _, o := range trace.OtherData {
			if o == nil {
				continue
			}
			if o.ID != trace.ID {
				continue
			}
			o.EndToEndProcessTime = timeToRun
		}
	}
	if c.options.uploadProfile {
		if err := trace.Upload(); err != nil {
			err = errors.Wrapf(err, "unable to upload profile file %s", profileFilePath)
			log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to upload profile output")
		}
	}
	return trace, nil
}

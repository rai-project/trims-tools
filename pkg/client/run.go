package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rai-project/uuid"

	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/assets"
	"github.com/rai-project/micro18-tools/pkg/config"
	"github.com/rai-project/micro18-tools/pkg/trace"
	"github.com/rai-project/micro18-tools/pkg/utils"
)

func Run(opts ...Option) ([]*trace.Trace, error) {
	hostname, _ := os.Hostname()

	options := WithOptions(opts...)

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
	for _, model := range models {
		var combined *trace.Trace
		dims, err := model.GetImageDimensions()
		if err != nil {
			dims = []uint32{3, 224, 224}
		}
		mean, err := model.GetMeanImage()
		if err != nil {
			mean = []float32{0, 0, 0}
		}
		if len(dims) != 3 {
			log.Errorf("expecting a 3 element vector for dimensions %v", dims)
			continue
		}
		cannonicalName := model.MustCanonicalName()
		progress.Prefix(fmt.Sprintf("running client model %s", cannonicalName))

		for ii := 0; ii < options.iterationCount; ii++ {

			progress.Increment()
			id := uuid.NewV4()
			profileFilePath := filepath.Join(config.Config.ProfileOutputDirectory, fmt.Sprintf("%s_%s.json", model.MustCanonicalName(), id))
			env := map[string]string{
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
				fmt.Sprintf("%s_%d", hostname, ii),
				strconv.Itoa(int(dims[0])),
				strconv.Itoa(int(dims[1])),
				strconv.Itoa(int(dims[2])),
				fmt.Sprintf("%v", mean[0]),
				fmt.Sprintf("%v", mean[1]),
				fmt.Sprintf("%v", mean[2]),
			)
			if !ran {
				log.WithError(errors.Errorf("failed to run cmd %s", config.Config.ClientRunCmd)).Error("failed to run model")
				continue
			}
			if err != nil {
				log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to run model")
				continue
			}
			if !options.postprocess {
				continue
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
			trace.Iteration = int64(ii)
			if options.uploadProfile {
				if err := trace.Upload(); err != nil {
					err = errors.Wrapf(err, "unable to upload profile file %s", profileFilePath)
					log.WithField("cmd", config.Config.ClientRunCmd).WithError(err).Error("failed to upload profile output")
				}
			}
			if combined == nil {
				combined = &trace
				combined.ID = uuid.NewV4()
			} else {
				combined.Combine(trace)
			}
		}
		if combined != nil {
			if err := combined.Upload(); err != nil {
				log.WithError(err).Error("failed to upload combined profile output")
			}
			res = append(res, combined)
		}
	}
	return res, nil
}

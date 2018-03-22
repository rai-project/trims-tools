package trace

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"strings"
	"time"

	"github.com/fatih/structs"

	"github.com/Unknwon/com"

	"hash/fnv"

	"github.com/Workiva/go-datastructures/augmentedtree"
	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/config"
	"github.com/rai-project/uuid"
	"github.com/spf13/cast"
	"github.com/ulule/deepcopier"
)

var (
	initTime               time.Time
	DefaultDisplayTimeUnit = "ms"
)

func timeUnit(unit string) (time.Duration, error) {
	switch unit {
	case "ns":
		return time.Nanosecond, nil
	case "us":
		return time.Microsecond, nil
	case "ms":
		return time.Millisecond, nil
	case "":
		return time.Microsecond, nil
	default:
		return time.Duration(0), errors.Errorf("the display time unit %v is not valid", unit)
	}
}

func mustTimeUnit(u string) time.Duration {
	unit, err := timeUnit(u)
	if err != nil {
		panic(err)
	}
	return unit
}

// Trace is an entry of trace format.
// https://github.com/catapult-project/catapult/tree/master/tracing
type TraceEvent struct {
	Name       string        `json:"name,omitempty"`
	Category   string        `json:"cat,omitempty"`
	EventType  string        `json:"ph,omitempty"`
	Timestamp  int64         `json:"ts"`  // displayTimeUnit
	Duration   time.Duration `json:"dur"` // displayTimeUnit
	ProcessID  int64         `json:"pid"`
	ThreadID   int64         `json:"tid,omitempty"`
	Args       interface{}   `json:"args,omitempty"`
	Stack      int           `json:"sf,omitempty"`
	EndStack   int           `json:"esf,omitempty"`
	Init       string        `json:"init_time,omitempty"`
	Start      int64         `json:"start,omitempty"`
	End        int64         `json:"end,omitempty"`
	InitTime   time.Time     `json:"init_time_t,omitempty"`
	StartTime  time.Time     `json:"start_time_t,omitempty"`
	EndTime    time.Time     `json:"end_time_t,omitempty"`
	Time       time.Time     `json:"time_t,omitempty"`
	TimeUnit   time.Duration `json:"timeUnit,omitempty"`
	UPREnabled bool          `json:"upr_enabled"`
	TraceID    string        `json:"trace_id,omitempty"`
}

type JSONTraceEvent struct {
	Name       string      `json:"name,omitempty"`
	Category   string      `json:"cat,omitempty"`
	EventType  string      `json:"ph,omitempty"`
	Timestamp  int64       `json:"ts,omitempty"`  // displayTimeUnit
	Duration   int64       `json:"dur,omitempty"` // displayTimeUnit
	ProcessID  int64       `json:"pid"`
	ThreadID   int64       `json:"tid,omitempty"`
	Args       interface{} `json:"args,omitempty"`
	Stack      int         `json:"sf,omitempty"`
	EndStack   int         `json:"esf,omitempty"`
	Init       string      `json:"init_time,omitempty"`
	Start      int64       `json:"start,omitempty"`
	End        int64       `json:"end,omitempty"`
	UPREnabled bool        `json:"upr_enabled,omitempty"`
}

type EventFrame struct {
	Name   string `json:"name"`
	Parent int    `json:"parent,omitempty"`
}

type TraceEvents []TraceEvent

type TraceServerInfo struct {
	ID               string    `json:"server_id,omitempty"`
	StartTime        time.Time `json:"server_start_time,omitempty"`
	EvictionPolicty  string    `json:"eviction_policy,omitempty"`
	EstimationRate   float32   `json:"estimation_rate,omitempty"`
	MemoryPercentage float32   `json:"memory_percentage,omitempty"`
	PersistCPU       bool      `json:"persist_cpu"`
}

type GitInfo struct {
	Commit string `json:"commit"`
	Date   string `json:"date"`
}

type TraceSummaryEvent struct {
	Name      string        `json:"name,omitempty"`
	Category  string        `json:"category,omitempty"`
	StartTime time.Time     `json:"start_time_t,omitempty"`
	Duration  time.Duration `json:"duration,omitempty"` // displayTimeUnit
}

// "load_params", "create", "set_input", "forward", "get_output_shape", "get_output", "free"
type TraceSummaryEvents struct {
	CreatePrediction      *TraceSummaryEvent `json:"create_prediction,omitempty"`
	LoadModel             *TraceSummaryEvent `json:"load_model,omitempty"`
	GRPCOpen              *TraceSummaryEvent `json:"grpc_open,omitempty"`
	NDArrayConvert        *TraceSummaryEvent `json:"ndarray_convert,omitempty"`
	Predict               *TraceSummaryEvent `json:"predict,omitempty"`
	LoadParams            *TraceSummaryEvent `json:"load_params,omitempty"`
	ReadParams            *TraceSummaryEvent `json:"read_params,omitempty"`
	SetPredictInput       *TraceSummaryEvent `json:"set_predict_input,omitempty"`
	ForwardPredict        *TraceSummaryEvent `json:"forward_predict,omitempty"`
	GetPredictOutputShape *TraceSummaryEvent `json:"get_output_shape,omitempty"`
	GetPredictOutput      *TraceSummaryEvent `json:"get_output,omitempty"`
	FreePredictor         *TraceSummaryEvent `json:"free_predictor,omitempty"`
}

type TraceSummary struct {
	ID                    string              `json:"run_id,omitempty"`
	ServerInfo            TraceServerInfo     `json:"server,omitempty"`
	EndToEndProcessTime   time.Duration       `json:"end_to_end_process_time,omitempty"`
	EndToEndTime          time.Duration       `json:"end_to_end_time,omitempty"`
	UPREnabled            bool                `json:"upr_enabled"`
	EagerMode             bool                `json:"eager_mode"`
	EagerModeAsync        bool                `json:"eager_mode_async"`
	Hostname              string              `json:"hostname"`
	ModelName             string              `json:"model_name"`
	ExperimentDescription string              `json:"experiment_description,omitempty"`
	InitTime              time.Time           `json:"init_time,omitempty"`
	Events                *TraceSummaryEvents `json:"events,omitempty"`
}

type TraceOtherData struct {
	ID                    string                 `json:"run_id,omitempty"`
	ServerInfo            TraceServerInfo        `json:"server,omitempty"`
	EndToEndProcessTime   time.Duration          `json:"end_to_end_process_time,omitempty"`
	EndToEndTime          time.Duration          `json:"end_to_end_time,omitempty"`
	UPREnabled            bool                   `json:"upr_enabled"`
	UPRBaseDirectory      string                 `json:"UPR_BASE_DIR"`
	EagerMode             bool                   `json:"eager_mode"`
	EagerModeAsync        bool                   `json:"eager_mode_async"`
	EndAt                 string                 `json:"end_at"`
	Git                   GitInfo                `json:"git"`
	Hostname              string                 `json:"hostname"`
	IsClient              bool                   `json:"is_client"`
	ModelName             string                 `json:"model_name"`
	ModelParams           string                 `json:"model_params"`
	ModelPath             string                 `json:"model_path"`
	StartAt               string                 `json:"start_at"`
	SymbolParams          string                 `json:"symbol_params"`
	Username              string                 `json:"username"`
	MinEvent              TraceEvent             `json:"min_event"`
	MaxEvent              TraceEvent             `json:"max_event"`
	Iteration             int64                  `json:"iteration,omitempty"`
	Input                 map[string]interface{} `json:"input,omitempty"`
	ExperimentDescription string                 `json:"experiment_description,omitempty"`
	IsProcessed           bool                   `json:"is_processed,omitempty"`
}

type Trace struct {
	ID                    string                `json:"id,omitempty"`
	ExperimentDescription string                `json:"experiment_description,omitempty"`
	UPREnabled            bool                  `json:"upr_enabled"`
	Iteration             int64                 `json:"iteration,omitempty"`
	StartTime             time.Time             `json:"start_time,omitempty"`
	EndTime               time.Time             `json:"end_time,omitempty"`
	TraceEvents           TraceEvents           `json:"traceEvents,omitempty"`
	DisplayTimeUnit       string                `json:"displayTimeUnit,omitempty"`
	Frames                map[string]EventFrame `json:"stackFrames"`
	TimeUnit              string                `json:"timeUnit,omitempty"`
	OtherDataRaw          *TraceOtherData       `json:"otherData,omitempty"`
	OtherData             []*TraceOtherData     `json:"otherDatas,omitempty"`
}

type JSONTrace struct {
	ID              string                `json:"id,omitempty"`
	UPREnabled      bool                  `json:"upr_enabled,omitempty"`
	TraceEvents     TraceEvents           `json:"traceEvents,omitempty"`
	DisplayTimeUnit string                `json:"displayTimeUnit,omitempty"`
	Frames          map[string]EventFrame `json:"stackFrames"`
	TimeUnit        string                `json:"timeUnit,omitempty"`
	OtherDataRaw    TraceOtherData        `json:"otherData,omitempty"`
	OtherData       []*TraceOtherData     `json:"otherDatas,omitempty"`
}

func (x *TraceEvent) UnmarshalJSON(data []byte) error {
	var jsonTraceEvent JSONTraceEvent
	err := json.Unmarshal(data, &jsonTraceEvent)
	if err != nil {
		log.WithError(err).Error("failed to unmarshal trace event data")
	}
	if err := deepcopier.Copy(jsonTraceEvent).To(x); err != nil {
		return errors.Wrapf(err, "unable to copy model")
	}
	x.Args = jsonTraceEvent.Args
	if x.Init == "" {
		return nil
	}
	timeUnit, err := timeUnit(DefaultDisplayTimeUnit)
	if err != nil {
		return err
	}
	x.TimeUnit = timeUnit
	initTime, err := time.Parse(time.RFC3339Nano, x.Init)
	if err != nil {
		return errors.Wrapf(err, "cannot parse time duration %s", x.Init)
	}
	x.InitTime = initTime
	x.Time = initTime.Add(time.Duration(x.Timestamp) * timeUnit)
	x.StartTime = initTime.Add(time.Duration(x.Start) * timeUnit)
	x.EndTime = initTime.Add(time.Duration(x.End) * timeUnit)
	if x.Duration == 0 {
		x.Duration = x.EndTime.Sub(x.StartTime)
	}
	return nil
}

// LowAtDimension returns an integer representing the lower bound
// at the requested dimension.
func (x TraceEvent) LowAtDimension(d uint64) int64 {
	if d != 1 {
		return 0
	}
	return x.Start
}

// HighAtDimension returns an integer representing the higher bound
// at the requested dimension.
func (x TraceEvent) HighAtDimension(d uint64) int64 {
	if d != 1 {
		return 0
	}
	return x.End
}

// OverlapsAtDimension should return a bool indicating if the provided
// interval overlaps this interval at the dimension requested.
func (mi TraceEvent) OverlapsAtDimension(iv augmentedtree.Interval, dimension uint64) bool {
	return mi.HighAtDimension(dimension) > iv.LowAtDimension(dimension) &&
		mi.LowAtDimension(dimension) < iv.HighAtDimension(dimension)
}

// ID should be a unique ID representing this interval.  This
// is used to identify which interval to delete from the tree if
// there are duplicates.
func (x TraceEvent) ID() uint64 {
	h := fnv.New64a()
	h.Write([]byte(x.Name))
	return h.Sum64()
}

func (t TraceEvents) Len() int           { return len(t) }
func (t TraceEvents) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TraceEvents) Less(i, j int) bool { return t[i].Timestamp < t[j].Timestamp }

func (t Trace) Len() int           { return t.TraceEvents.Len() }
func (t Trace) Swap(i, j int)      { t.TraceEvents.Swap(i, j) }
func (t Trace) Less(i, j int) bool { return t.TraceEvents.Less(i, j) }

func (x *Trace) UnmarshalJSON(data []byte) error {
	var jsonTrace JSONTrace

	err := json.Unmarshal(data, &jsonTrace)
	if err != nil {
		log.WithError(err).Error("failed to unmarshal trace data")
	}
	if err := deepcopier.Copy(jsonTrace).To(x); err != nil {
		return errors.Wrapf(err, "unable to copy model")
	}
	id := jsonTrace.ID
	if id == "" {
		id = jsonTrace.OtherDataRaw.ID
	}
	if id == "" {
		id = uuid.NewV4()
	}
	experimentDescription := config.Config.ExperimentDescription

	x.ID = id
	x.ExperimentDescription = experimentDescription
	for ii := range x.TraceEvents {
		x.TraceEvents[ii].TraceID = id
	}
	if len(jsonTrace.OtherData) != 0 {
		x.OtherData = make([]*TraceOtherData, len(jsonTrace.OtherData))
		for ii := range jsonTrace.OtherData {
			x.OtherData[ii] = new(TraceOtherData)
			if err := deepcopier.Copy(jsonTrace.OtherData[ii]).To(x.OtherData[ii]); err != nil {
				return errors.Wrapf(err, "unable to copy other data model %d", ii)
			}
		}
		return nil
	}
	x.OtherDataRaw = new(TraceOtherData)
	if err := deepcopier.Copy(jsonTrace.OtherDataRaw).To(x.OtherDataRaw); err != nil {
		return errors.Wrapf(err, "unable to copy other data model")
	}
	if x.OtherDataRaw.IsProcessed {
		return nil
	}
	x.OtherDataRaw.IsProcessed = true
	if !jsonTrace.OtherDataRaw.IsClient {
		x.OtherDataRaw.ServerInfo = TraceServerInfo{}
	}
	if jsonTrace.OtherDataRaw.IsClient && config.Config.UPREnabled {
		serverInfoPath := config.Config.ServerInfoPath
		if com.IsFile(serverInfoPath) {
			bts, err := ioutil.ReadFile(serverInfoPath)
			if err == nil {
				var info TraceServerInfo
				if err := json.Unmarshal(bts, &info); err == nil {
					if x.StartTime.Sub(info.StartTime) > 0 {
						x.OtherDataRaw.ServerInfo = info
					}
				}
			}
		}
	}
	x.OtherDataRaw.ID = id
	x.OtherDataRaw.ExperimentDescription = experimentDescription
	x.StartTime, _ = time.Parse(time.RFC3339Nano, x.OtherDataRaw.StartAt)
	x.EndTime, _ = time.Parse(time.RFC3339Nano, x.OtherDataRaw.EndAt)

	minEvent := x.MinEvent()
	maxEvent := x.MaxEvent()
	if x.OtherDataRaw.EndToEndTime == 0 {
		x.OtherDataRaw.EndToEndTime = maxEvent.EndTime.Sub(minEvent.StartTime)
	}
	if x.OtherDataRaw.EndToEndProcessTime == 0 {
		x.OtherDataRaw.EndToEndProcessTime = x.EndTime.Sub(x.StartTime)
	}
	if structs.IsZero(x.OtherDataRaw.MaxEvent) {
		x.OtherDataRaw.MaxEvent = maxEvent
	}
	if structs.IsZero(x.OtherDataRaw.MinEvent) {
		x.OtherDataRaw.MinEvent = minEvent
	}

	x.OtherData = []*TraceOtherData{x.OtherDataRaw}

	return nil
}

func (x Trace) Summarize(isFull bool) (*TraceSummary, error) {
	data := x.OtherDataRaw
	if data == nil {
		return nil, errors.New("expecting OtherDataRaw to be in the trace. summary does not work on combined output")
	}
	summary, err := data.Summarize()
	if err != nil {
		return nil, err
	}
	initTime, err := time.Parse(time.RFC3339Nano, data.StartAt)
	if err == nil {
		summary.InitTime = initTime
	}
	if isFull {
		summary.fillEvents(x)
	}
	return summary, nil
}

func (x TraceOtherData) Summarize() (*TraceSummary, error) {
	summary := new(TraceSummary)
	if err := deepcopier.Copy(x).To(summary); err != nil {
		return nil, errors.Wrap(err, "unable to copy summarize trace")
	}
	return summary, nil
}

func (x TraceEvent) Summarize() *TraceSummaryEvent {
	return &TraceSummaryEvent{
		Name:      x.Name,
		Category:  x.Category,
		StartTime: x.StartTime,
		Duration:  x.Duration,
	}
}

func (x *TraceSummary) fillEvents(tr Trace) {
	if x.Events == nil {
		x.Events = new(TraceSummaryEvents)
	}
	se := x.Events
	events := tr.TraceEvents
	for _, event := range events {
		if event.EventType != "B" {
			continue
		}
		switch event.Category {
		case "prediction":
			switch event.Name {
			case "load_params":
				se.LoadParams = event.Summarize()
				continue
			case "create":
				se.CreatePrediction = event.Summarize()
				continue
			case "set_input":
				se.SetPredictInput = event.Summarize()
				continue
			case "forward":
				se.ForwardPredict = event.Summarize()
				continue
			case "predict":
				se.Predict = event.Summarize()
				continue
			case "get_output_shape":
				se.GetPredictOutputShape = event.Summarize()
				continue
			case "get_output":
				se.GetPredictOutput = event.Summarize()
				continue
			case "free":
				se.FreePredictor = event.Summarize()
				continue
			default:
				continue
			}
		case "serialization":
			switch event.Name {
			case "convering_to_nd_array":
				se.NDArrayConvert = event.Summarize()
				continue
			default:
				continue
			}
		case "grpc":
			switch event.Name {
			case "open":
				se.GRPCOpen = event.Summarize()
				continue
			default:
				continue
			}
		case "io":
			switch event.Name {
			case "read_params":
				se.ReadParams = event.Summarize()
				continue
			default:
				continue
			}
		case "load":
			switch event.Name {
			case "load_model":
				se.LoadModel = event.Summarize()
				continue
			default:
				continue
			}
		default:
			// ignore
		}
	}
}

func (x Trace) Adjust() (Trace, error) {
	tr, err := x.DeleteIgnoredEvents()
	if err != nil {
		log.WithError(err).Error("failed to delete ignored events")
		tr = x
	}
	return tr.UpdateEventNames().ZeroOut(), nil
}

func (x Trace) DeleteIgnoredEvents() (Trace, error) {
	var minTimeStamp int64
	var adjustedEvent TraceEvent
	events := TraceEvents{}
	for _, event := range x.TraceEvents { // assumes that there is only one thing to ignore
		if event.Category == "ignore" {
			if event.EventType == "E" {
				adjustedEvent = event
			}
		}
		if minTimeStamp < event.Timestamp {
			minTimeStamp = event.Timestamp
		}
	}
	if adjustedEvent.Name == "" {
		return x, nil
	}
	// pp.Println(adjustedEvent)
	for _, event := range x.TraceEvents {
		timeUnit := event.TimeUnit
		// initTime, _ := time.Parse(time.RFC3339Nano, event.Init)
		// pp.Println(timeAdjustmentI, "   ", event.Timestamp, "   ", event.Timestamp-timeAdjustmentI)
		if event.Category == "ignore" {
			continue
		}
		if event.EndTime.After(adjustedEvent.Time) && event.StartTime.Before(adjustedEvent.Time) {
			event.Duration = event.Duration - adjustedEvent.Duration
		}
		if event.EndTime.Before(adjustedEvent.Time) {
			events = append(events, event)
			continue
		}
		// if event.Name == "load_nd_array" {
		// 	continue
		// }

		if event.EventType == "B" || event.EventType == "E" {
			if event.Time.After(adjustedEvent.Time) {
				event.Time = event.Time.Add(-adjustedEvent.Duration)
				event.Timestamp = event.Timestamp - int64(adjustedEvent.Duration/timeUnit)
				// pp.Println(event.Timestamp, "   ", adjustedEvent.Timestamp, "  ", int64(adjustedEvent.Duration), "   ", event.Timestamp-adjustedEvent.Timestamp+minTimeStamp)
			}
			if event.StartTime.After(adjustedEvent.StartTime) {
				event.Start = event.Start - adjustedEvent.Start
				event.StartTime = event.StartTime.Add(-adjustedEvent.Duration)
			}

			if event.EndTime.After(adjustedEvent.EndTime) {
				event.End = event.End - adjustedEvent.End
				event.EndTime = event.EndTime.Add(-adjustedEvent.Duration)
			}
		}

		events = append(events, event)
	}
	x.TraceEvents = events
	return x, nil
}

func (x Trace) MaxEvent() TraceEvent {
	var maxEvent TraceEvent
	maxTimeStamp := int64(math.MinInt64)
	for _, event := range x.TraceEvents { // assumes that there is only one thing to ignore
		if event.Category == "ignore" {
			continue
		}
		if event.EventType != "E" {
			continue
		}
		if maxTimeStamp < event.Timestamp {
			maxTimeStamp = event.Timestamp
			maxEvent = event
		}
	}
	return maxEvent
}

func (x Trace) MinEvent() TraceEvent {
	var minEvent TraceEvent
	minTimeStamp := int64(math.MaxInt64)
	for _, event := range x.TraceEvents { // assumes that there is only one thing to ignore
		if event.Category == "ignore" {
			continue
		}
		if event.EventType != "B" {
			continue
		}
		if minTimeStamp > event.Timestamp {
			minTimeStamp = event.Timestamp
			minEvent = event
		}
	}
	return minEvent
}

func (x Trace) ZeroOut() Trace {
	minEvent := x.MinEvent()
	minTimeStamp := minEvent.Timestamp

	td := minEvent.StartTime.Sub(x.StartTime)

	return x.AddTimestampOffset(-minTimeStamp).AddDurationOffset(td)
}

func (x Trace) AddTimestampOffset(ts int64) Trace {
	events := make([]TraceEvent, len(x.TraceEvents))
	for ii, event := range x.TraceEvents {
		if event.EventType == "B" || event.EventType == "E" {
			event.Timestamp = event.Timestamp + ts
			event.Start = event.Start + ts
			event.End = event.End + ts
		}
		events[ii] = event
	}
	x.TraceEvents = events
	return x
}

func (x Trace) AddDurationOffset(td time.Duration) Trace {
	events := make([]TraceEvent, len(x.TraceEvents))
	for ii, event := range x.TraceEvents {
		if event.EventType == "B" || event.EventType == "E" {
			event.Time = event.Time.Add(td)
			event.StartTime = event.StartTime.Add(td)
			event.EndTime = event.EndTime.Add(td)
		}
		events[ii] = event
	}
	x.TraceEvents = events
	return x
}

func (tr Trace) UpdateEventNames() Trace {
	events := make([]TraceEvent, len(tr.TraceEvents))
	for ii, event := range tr.TraceEvents {
		if event.EventType == "M" {
			name := tr.ID
			otherData := tr.OtherDataRaw
			if otherData == nil || otherData.ModelName == "" {
				for _, e := range tr.OtherData {
					if e.ID == e.ID {
						otherData = e
						break
					}
				}
			}
			if otherData != nil {
				uprEnabled := "upr_enabled=" + cast.ToString(otherData.UPREnabled)
				modelName := "model_name=" + otherData.ModelName
				hostName := "host_name=" + otherData.Hostname
				iteration := "iteration=" + cast.ToString(otherData.Iteration)
				name = strings.Join([]string{modelName, hostName, uprEnabled, iteration}, ",")
			}
			event.Args = map[string]string{
				"name":        name,
				"upr_enabled": cast.ToString(tr.UPREnabled),
			}
		}
		events[ii] = event
	}
	tr.TraceEvents = events
	return tr
}

func (x Trace) HashID() int64 {
	if x.Iteration != 0 {
		return x.Iteration
	}
	h := fnv.New32a()
	h.Write([]byte(x.ID))
	return int64(h.Sum32())
}

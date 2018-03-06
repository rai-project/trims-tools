package trace

import (
	"encoding/json"
	"time"

	"hash/fnv"

	"github.com/Workiva/go-datastructures/augmentedtree"
	"github.com/pkg/errors"
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

// Trace is an entry of trace format.
// https://github.com/catapult-project/catapult/tree/master/tracing
type TraceEvent struct {
	Name      string            `json:"name,omitempty"`
	Category  string            `json:"cat,omitempty"`
	EventType string            `json:"ph,omitempty"`
	Timestamp int64             `json:"ts,omitempty"`  // displayTimeUnit
	Duration  time.Duration     `json:"dur,omitempty"` // displayTimeUnit
	ProcessID int64             `json:"pid,omitempty"`
	ThreadID  int64             `json:"tid,omitempty"`
	Args      map[string]string `json:"args,omitempty"`
	Stack     int               `json:"sf,omitempty"`
	EndStack  int               `json:"esf,omitempty"`
	Init      string            `json:"init_time,omitempty"`
	Start     int64             `json:"begin,omitempty"`
	End       int64             `json:"end,omitempty"`
	InitTime  time.Time         `json:"-"`
	StartTime time.Time         `json:"-"`
	EndTime   time.Time         `json:"-"`
	Time      time.Time         `json:"-"`
	TimeUnit  time.Duration     `json:"-"`
}

func (x *TraceEvent) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, x)
	if err != nil {
		return err
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

type EventFrame struct {
	Name   string `json:"name"`
	Parent int    `json:"parent,omitempty"`
}

type TraceEvents []TraceEvent

func (t TraceEvents) Len() int           { return len(t) }
func (t TraceEvents) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TraceEvents) Less(i, j int) bool { return t[i].Timestamp < t[j].Timestamp }

type TraceOtherData struct {
	UPRBaseDirectory string `json:"UPR_BASE_DIR"`
	EagerMode        bool   `json:"eager_mode"`
	EagerModeAsync   bool   `json:"eager_mode_async"`
	EndAt            string `json:"end_at"`
	Git              struct {
		Commit string `json:"commit"`
		Date   string `json:"date"`
	} `json:"git"`
	Hostname     string `json:"hostname"`
	IsClient     bool   `json:"is_client"`
	ModelName    string `json:"model_name"`
	ModelParams  string `json:"model_params"`
	ModelPath    string `json:"model_path"`
	StartAt      string `json:"start_at"`
	SymbolParams string `json:"symbol_params"`
	Username     string `json:"username"`
}

//easyjson:json
type Trace struct {
	ID              string                `json:"-"`
	StartTime       time.Time             `json:"-"`
	EndTime         time.Time             `json:"-"`
	TraceEvents     TraceEvents           `json:"traceEvents,omitempty"`
	DisplayTimeUnit string                `json:"displayTimeUnit,omitempty"`
	Frames          map[string]EventFrame `json:"stackFrames"`
	TimeUnit        string                `json:"timeUnit,omitempty"`
	OtherData       TraceOtherData        `json:"otherData,omitempty"`
}

func (t Trace) Len() int           { return t.TraceEvents.Len() }
func (t Trace) Swap(i, j int)      { t.TraceEvents.Swap(i, j) }
func (t Trace) Less(i, j int) bool { return t.TraceEvents.Less(i, j) }

func (x *Trace) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, x)
	if err != nil {
		return err
	}
	return nil
}

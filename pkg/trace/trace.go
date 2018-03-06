//go:generate jsonenums -type=TraceEvent
//go:generate jsonenums -type=EventFrame
//go:generate jsonenums -type=TraceOtherData
//go:generate jsonenums -type=Trace
package trace

import (
	"fmt"
	"time"
)

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
	Start     int64             `json:"begin,omitempty"`
	End       int64             `json:"end,omitempty"`
	StartTime time.Time         `json:"-"`
	EndTime   time.Time         `json:"-"`
	Time      time.Time         `json:"-"`
}

type EventFrame struct {
	Name   string `json:"name"`
	Parent int    `json:"parent,omitempty"`
}

func (t TraceEvent) ID() string {
	return fmt.Sprintf("%s::%s/%v", t.Category, t.Name, t.ThreadID)
}

type TraceEvents []TraceEvent

func (t TraceEvents) Len() int           { return len(t) }
func (t TraceEvents) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TraceEvents) Less(i, j int) bool { return t[i].Timestamp < t[j].Timestamp }

type TraceOtherData struct {
	UPRBASEDIR     string `json:"UPR_BASE_DIR"`
	EagerMode      bool   `json:"eager_mode"`
	EagerModeAsync bool   `json:"eager_mode_async"`
	EndAt          string `json:"end_at"`
	Git            struct {
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

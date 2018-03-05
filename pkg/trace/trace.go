package trace

import (
	"fmt"
	"time"
)

// Trace is an entry of trace format.
// https://github.com/catapult-project/catapult/tree/master/tracing
//easyjson:json
type TraceEvent struct {
	Name      string                 `json:"name,omitempty"`
	Category  string                 `json:"cat,omitempty"`
	EventType string                 `json:"ph,omitempty"`
	Timestamp int64                  `json:"ts,omitempty"`  // displayTimeUnit
	Duration  time.Duration          `json:"dur,omitempty"` // displayTimeUnit
	ProcessID int64                  `json:"pid,omitempty"`
	ThreadID  int64                  `json:"tid,omitempty"`
	Args      map[string]interface{} `json:"args,omitempty"`
	Stack     int                    `json:"sf,omitempty"`
	EndStack  int                    `json:"esf,omitempty"`
	Time      time.Time              `json:"-"`
}

//easyjson:json
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

//easyjson:json
type Trace struct {
	StartTime       time.Time              `json:"-"`
	EndTime         time.Time              `json:"-"`
	TraceEvents     TraceEvents            `json:"traceEvents,omitempty"`
	DisplayTimeUnit string                 `json:"displayTimeUnit,omitempty"`
	Frames          map[string]EventFrame  `json:"stackFrames"`
	TimeUnit        string                 `json:"timeUnit,omitempty"`
	OtherData       map[string]interface{} `json:"otherData,omitempty"`
}

func (t Trace) Len() int           { return t.TraceEvents.Len() }
func (t Trace) Swap(i, j int)      { t.TraceEvents.Swap(i, j) }
func (t Trace) Less(i, j int) bool { return t.TraceEvents.Less(i, j) }

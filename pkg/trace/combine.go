package trace

func Combine(trace0 Trace, others ...Trace) *Trace {
	trace := &trace0
	for _, o := range others {
		trace.Combine(o)
	}
	return trace
}

func (tr *Trace) Combine(other Trace) {
	startTime := tr.StartTime
	tr.OtherData = append(tr.OtherData, other.OtherData...)
	timeUnit := mustTimeUnit(tr.TimeUnit)
	for _, event := range other.TraceEvents {
		event.ProcessID = int64(other.HashID())
		// event.InitTime = startTime
		// event.StartTime = initTime.Add(time.Duration(event.Start) * timeUnit)
		// event.EndTime = initTime.Add(time.Duration(event.End) * timeUnit)
		tr.TraceEvents = append(tr.TraceEvents, event)
	}
	_ = startTime
	_ = timeUnit
	tr.OtherDataRaw = nil
}

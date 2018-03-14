package trace

func Combine(trace0 Trace, others ...Trace) *Trace {
	trace := &trace0
	for _, o := range others {
		trace.Combine(o)
	}
	return trace
}

func (tr *Trace) Combine(other Trace) {
	for _, event := range other.TraceEvents {
		event.ProcessID = int64(len(tr.OtherData))
		tr.TraceEvents = append(tr.TraceEvents, event)
	}
	tr.OtherData = append(tr.OtherData, other.OtherData...)
	tr.OtherDataRaw = nil
}

package metrics

type Metrics interface {
	IncEvent(eventType string)
	IncSessionStart()
	IncSessionComplete()
	IncStateTransition(prev, next string)
	IncError(errorType string)
	Snapshot() Snapshot
}

type Snapshot struct {
	TotalEvents      int64
	EventsByType     map[string]int64
	SessionsStarted  int64
	SessionsComplete int64
	StateTransitions map[string]int64
	ErrorsByType     map[string]int64
}

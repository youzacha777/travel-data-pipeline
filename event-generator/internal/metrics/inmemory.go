package metrics

import (
	"sync"
	"sync/atomic"
)

type InMemoryMetrics struct {
	totalEvents      atomic.Int64
	sessionsStarted  atomic.Int64
	sessionsComplete atomic.Int64

	// key: string, value: *atomic.Int64
	eventsByType     sync.Map
	stateTransitions sync.Map
}

func NewInMemory() *InMemoryMetrics {
	return &InMemoryMetrics{}
}

// =======================
// Increment methods
// =======================

func (m *InMemoryMetrics) IncEvent(eventType string) {
	m.totalEvents.Add(1)

	val, _ := m.eventsByType.LoadOrStore(eventType, &atomic.Int64{})
	val.(*atomic.Int64).Add(1)

	// 이벤트 발생 시 로그 출력 (확인용)
	// fmt.Printf("Event of type '%s' incremented\n", eventType)
}

func (m *InMemoryMetrics) IncSessionStart() {
	m.sessionsStarted.Add(1)
}

func (m *InMemoryMetrics) IncSessionComplete() {
	m.sessionsComplete.Add(1)
}

func (m *InMemoryMetrics) IncStateTransition(prev, next string) {
	key := prev + " -> " + next

	val, _ := m.stateTransitions.LoadOrStore(key, &atomic.Int64{})
	val.(*atomic.Int64).Add(1)
}

// =======================
// Snapshot
// =======================

func (m *InMemoryMetrics) Snapshot() Snapshot {
	snap := Snapshot{
		TotalEvents:      m.totalEvents.Load(),
		SessionsStarted:  m.sessionsStarted.Load(),
		SessionsComplete: m.sessionsComplete.Load(),
		EventsByType:     make(map[string]int64),
		StateTransitions: make(map[string]int64),
	}

	m.eventsByType.Range(func(k, v any) bool {
		snap.EventsByType[k.(string)] = v.(*atomic.Int64).Load()
		return true
	})

	m.stateTransitions.Range(func(k, v any) bool {
		snap.StateTransitions[k.(string)] = v.(*atomic.Int64).Load()
		return true
	})

	return snap
}

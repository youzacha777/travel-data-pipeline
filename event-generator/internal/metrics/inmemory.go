package metrics

import (
	"sync"
	"sync/atomic"
)

// InMemoryMetrics 는 Metrics 인터페이스를 구현하며
// 이벤트, 세션, 상태 전환, 에러 카운트를 기록합니다.
type InMemoryMetrics struct {
	totalEvents      atomic.Int64
	sessionsStarted  atomic.Int64
	sessionsComplete atomic.Int64

	eventsByType     sync.Map
	stateTransitions sync.Map
	errorsByType     sync.Map
}

// NewInMemory 초기화
func NewInMemory() *InMemoryMetrics {
	return &InMemoryMetrics{}
}

// =======================
// Increment methods
// =======================

// 이벤트 카운트
func (m *InMemoryMetrics) IncEvent(eventType string) {
	m.totalEvents.Add(1)
	val, _ := m.eventsByType.LoadOrStore(eventType, &atomic.Int64{})
	val.(*atomic.Int64).Add(1)
}

// 세션 시작 카운트
func (m *InMemoryMetrics) IncSessionStart() {
	m.sessionsStarted.Add(1)
}

// 세션 완료 카운트
func (m *InMemoryMetrics) IncSessionComplete() {
	m.sessionsComplete.Add(1)
}

// 상태 전환 카운트
func (m *InMemoryMetrics) IncStateTransition(prev, next string) {
	key := prev + " -> " + next
	val, _ := m.stateTransitions.LoadOrStore(key, &atomic.Int64{})
	val.(*atomic.Int64).Add(1)
}

// 에러 카운트
func (m *InMemoryMetrics) IncError(errorType string) {
	val, _ := m.errorsByType.LoadOrStore(errorType, &atomic.Int64{})
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
		ErrorsByType:     make(map[string]int64),
	}

	m.eventsByType.Range(func(k, v any) bool {
		snap.EventsByType[k.(string)] = v.(*atomic.Int64).Load()
		return true
	})

	m.stateTransitions.Range(func(k, v any) bool {
		snap.StateTransitions[k.(string)] = v.(*atomic.Int64).Load()
		return true
	})

	m.errorsByType.Range(func(k, v any) bool {
		snap.ErrorsByType[k.(string)] = v.(*atomic.Int64).Load()
		return true
	})

	return snap
}

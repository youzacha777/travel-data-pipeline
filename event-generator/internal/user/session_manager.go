package user

import (
	"event-generator/internal/event"
	"event-generator/internal/fsm"
	"event-generator/internal/metrics"
	"sync"
	"time"
)

// =======================
// PayloadGenerator Interface
// =======================
type PayloadGenerator interface {
	Generate(eventType string, session *Session) map[string]any
}

// =======================
// SessionManager
// =======================
type SessionManager struct {
	userPool   *UserPool
	fsm        fsm.FSM
	payloadGen PayloadGenerator
	eventChan  chan *event.Event
	metrics    metrics.Metrics

	sessions map[string]*Session // key: sessionID
	ttl      time.Duration

	mu sync.Mutex
}

// =======================
// Constructor
// =======================
func NewSessionManager(
	userPool *UserPool,
	fsm fsm.FSM,
	payloadGen PayloadGenerator,
	eventChan chan *event.Event,
	metricStore metrics.Metrics,
	ttl time.Duration,
) *SessionManager {
	return &SessionManager{
		userPool:   userPool,
		fsm:        fsm,
		payloadGen: payloadGen, // 인터페이스라 *generator.PayloadGenerator 가능
		eventChan:  eventChan,
		metrics:    metricStore,
		sessions:   make(map[string]*Session),
		ttl:        ttl,
	}
}

// =======================
// Public API
// =======================
func (sm *SessionManager) Step() {
	now := time.Now().UnixMilli()
	u := sm.userPool.GetRandomUser()
	if u == nil {
		return
	}

	s := sm.getOrCreateSession(u.ID, now)

	// FSM 상태 전이
	ev := sm.fsm.Step(s, now)
	if ev == nil {
		sm.cleanupSessions(now)
		return
	}

	// 상태 전환이 있을 때만 트랜지션 기록
	if ev.Attributes.PrevState != ev.Attributes.State {
		sm.metrics.IncStateTransition(ev.Attributes.PrevState, ev.Attributes.State)
	}

	// 세션 시작 시 메트릭 기록
	if ev.EventType == string(fsm.EventSearchSubmitted) {
		sm.metrics.IncSessionStart()
	}

	// [중요] payload 생성 시 세션과 이벤트 타입을 함께 전달
	// (Generate 함수가 내부에서 genPurchase(s, ev.EventType)를 호출하게 됩니다)
	payload := sm.payloadGen.Generate(ev.EventType, s)

	// 페이로드 합치기 로직...
	if ev.Attributes.Extra == nil {
		ev.Attributes.Extra = make(map[string]any)
	}
	for k, v := range payload {
		ev.Attributes.Extra[k] = v
	}

	sm.eventChan <- ev

	// [중요] 최종 종료 이벤트인 경우에만 여기서 세션을 명시적으로 삭제
	if ev.EventType == string(fsm.EventExit) {
		sm.mu.Lock()
		delete(sm.sessions, s.ID)
		sm.mu.Unlock()
		sm.metrics.IncSessionComplete() // 세션 종료 시 메트릭 기록
	}

	sm.cleanupSessions(now)
}

// =======================
// Internal helpers
// =======================

func (sm *SessionManager) getOrCreateSession(userID string, now int64) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 1. 기존 세션 중 '진짜 종료(StateExit)' 상태가 아닌 세션을 찾음
	for _, s := range sm.sessions {
		// StatePurchase 상태여도 데이터를 유지해야 하므로 재사용 대상으로 포함합니다.
		if s.UserID == userID && s.State != fsm.StateExit {
			s.LastEventTs = now
			s.ExpiresAt = now + sm.ttl.Milliseconds()
			return s
		}
	}

	// 2. StateExit 상태이거나 세션이 아예 없을 때만 새로 생성
	sessionID := generateSessionID()
	s := NewSession(sessionID, userID, sm.ttl)
	s.SetState(fsm.StateBrowsing)
	sm.sessions[sessionID] = s

	// 3. 새로운 세션이 시작될 때 메트릭을 증가시킴
	if sm.metrics != nil {
		sm.metrics.IncSessionStart()
	}

	return s
}

func (sm *SessionManager) cleanupSessions(now int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for sessionID, s := range sm.sessions {
		if s.ExpiresAt <= now {
			// 세션 종료 메트릭 기록
			if sm.metrics != nil {
				sm.metrics.IncSessionComplete()
			}
			delete(sm.sessions, sessionID)
		}
	}
}

// =======================
// Utils
// =======================
func generateSessionID() string {
	return time.Now().Format("20060102150405.000000000")
}

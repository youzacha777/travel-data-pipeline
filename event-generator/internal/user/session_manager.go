package user

import (
	"event-generator/internal/event"
	"event-generator/internal/fsm"
	"event-generator/internal/metrics"
	"fmt"
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

	// [수정] 맵 구조 개선: userID로 세션ID를 즉시 찾기 위한 인덱스 추가
	sessions      map[string]*Session // key: sessionID
	userToSession map[string]string   // key: userID, value: sessionID
	ttl           time.Duration

	mu sync.RWMutex // [수정] 읽기 성능 향상을 위해 RWMutex 사용
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
	sm := &SessionManager{
		userPool:      userPool,
		fsm:           fsm,
		payloadGen:    payloadGen,
		eventChan:     eventChan,
		metrics:       metricStore,
		sessions:      make(map[string]*Session),
		userToSession: make(map[string]string), // 맵 초기화
		ttl:           ttl,
	}

	// [수정] 매 Step마다 하던 청소를 별도 고루틴으로 분리 (워커 부하 감소)
	go sm.backgroundCleanup()

	return sm
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

	// 1. 세션 조회/생성 (O(1) 성능)
	s := sm.getOrCreateSession(u.ID, now)

	// 2. FSM 상태 전이
	ev := sm.fsm.Step(s, now)
	if ev == nil {
		// ev가 nil이면 세션이 종료되었거나 유효하지 않은 상태
		return
	}

	// 상태 전환 및 메트릭 기록
	if ev.Attributes.PrevState != ev.Attributes.State {
		sm.metrics.IncStateTransition(ev.Attributes.PrevState, ev.Attributes.State)
	}

	if ev.EventType == string(fsm.EventSearchSubmitted) {
		sm.metrics.IncSessionStart()
	}

	// 페이로드 생성 및 병합
	payload := sm.payloadGen.Generate(ev.EventType, s)
	if ev.Attributes.Extra == nil {
		ev.Attributes.Extra = make(map[string]any)
	}
	for k, v := range payload {
		ev.Attributes.Extra[k] = v
	}

	// 3. 채널 전송
	sm.eventChan <- ev

	// 4. 종료 이벤트인 경우 즉시 삭제
	if ev.EventType == string(fsm.EventExit) {
		sm.deleteSession(u.ID, s.ID)
		sm.metrics.IncSessionComplete()
	}

}

// =======================
// Internal helpers
// =======================

func (sm *SessionManager) getOrCreateSession(userID string, now int64) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// userID 인덱스를 통해 O(1)로 조회
	if sid, ok := sm.userToSession[userID]; ok {
		if s, exists := sm.sessions[sid]; exists && s.State != fsm.StateExit {
			s.LastEventTs = now
			s.ExpiresAt = now + sm.ttl.Milliseconds()
			return s
		}
	}

	// 기존 세션이 없으면 새로 생성
	sessionID := fmt.Sprintf("sess_%s_%d", userID, now)
	s := NewSession(sessionID, userID, sm.ttl)
	s.SetState(fsm.StateBrowsing)

	sm.sessions[sessionID] = s
	sm.userToSession[userID] = sessionID

	if sm.metrics != nil {
		sm.metrics.IncSessionStart()
	}

	return s
}

// 세션 명시적 삭제
func (sm *SessionManager) deleteSession(userID, sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, sessionID)
	delete(sm.userToSession, userID)
}

// 백그라운드 세션 청소 (워커들의 락 경합 방지)
func (sm *SessionManager) backgroundCleanup() {
	ticker := time.NewTicker(2 * time.Second) // 2초마다 수행
	for range ticker.C {
		now := time.Now().UnixMilli()
		sm.mu.Lock()
		for sid, s := range sm.sessions {
			if s.ExpiresAt <= now {
				delete(sm.userToSession, s.UserID)
				delete(sm.sessions, sid)
				if sm.metrics != nil {
					sm.metrics.IncSessionComplete()
				}
			}
		}
		sm.mu.Unlock()
	}
}

// ID 생성
func generateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

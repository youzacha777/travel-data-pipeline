package fsm

import (
	"event-generator/internal/event"
	"fmt"
	"math/rand/v2" // v1 대신 v2를 사용합니다. (Thread-safe 최적화)
)

// =======================================================
// Session Interface (FSM이 요구하는 최소 계약)
// =======================================================
type Session interface {
	GetID() string
	GetUserID() string

	GetState() State
	SetState(State)

	GetPrevState() State
	SetPrevState(State)

	GetLastEventTs() int64
	SetLastEventTs(int64)

	GetSearchKeyword() string
	SetSearchKeyword(string)

	SetEventPage(string)
	GetEventPage() string
	SetPageType(string)
	GetPageType() string
	SetBrowsingCountryCategory(string)
	GetBrowsingCountryCategory() string
	SetBrowsingProductCategory(string)
	GetBrowsingProductCategory() string
	ResetBrowsingContext()

	SetLastPicked(productID, category, country string)
	GetLastProductID() string
	GetLastCategory() string
	GetLastCountry() string
	GetLastPicked() (productID, category, country string)

	SetLastQuantity(qty int)
	GetLastQuantity() int

	GetPageIndex() int
	SetPageIndex(int)
	IncrementPageIndex()

	SetExpiresAt(int64)
}

// =======================================================
// FSM Interface
// =======================================================
type FSM interface {
	Step(s Session, now int64) *event.Event
}

// =======================================================
// SimpleFSM
// =======================================================
type SimpleFSM struct {
	// 이제 rnd *rand.Rand를 저장할 필요가 없습니다.
	// rand/v2의 전역 함수는 내부적으로 성능 최적화와 동시성 보호가 되어 있습니다.
}

func NewSimpleFSM() *SimpleFSM {
	return &SimpleFSM{}
}

// =======================================================
// Step
// =======================================================
func (f *SimpleFSM) Step(s Session, now int64) *event.Event {

	// 1. terminal state 처리
	transitions, ok := Transitions[s.GetState()]
	if !ok || len(transitions) == 0 {
		return nil
	}

	// 2. 최초 상태가 없다면 StateBrowsing으로 설정
	if s.GetState() == StateNone {
		s.SetState(StateBrowsing)
	}

	// 3. transition 선택
	// f.rnd 대신 전역 rand를 사용하도록 chooseTransition의 인자를 수정해야 합니다.
	tr := chooseTransition(transitions)
	if tr == nil {
		return nil
	}

	// 4. 이전 상태 보존
	prevState := s.GetState()
	nextState := tr.NextState
	evType := tr.Event

	// 5. Back 이벤트 처리
	if evType == EventBack {
		if s.GetPrevState() != StateNone {
			nextState = s.GetPrevState()
		} else {
			nextState = StateBrowsing
		}
	}

	// 6. 세션 상태 갱신
	s.SetPrevState(prevState)
	s.SetState(nextState)
	s.SetLastEventTs(now)

	// 7. 검색 이벤트일 때만 키워드 생성
	if evType == EventSearchSubmitted {
		// FakeKeyword 역시 내부적으로 rand/v2를 쓰도록 수정 필요합니다.
		s.SetSearchKeyword(FakeKeyword())
		s.SetPageIndex(1)
	}

	// 8. 이벤트 생성
	// rand.Int63n -> rand.Int64N (v2 대문자 N)으로 변경하여 안전하게 ID 생성
	return &event.Event{
		EventID:   fmt.Sprintf("evt-%d-%09d", now, rand.Int64N(1_000_000_000)),
		EventType: string(evType),
		EventTs:   now,
		UserID:    s.GetUserID(),
		SessionID: s.GetID(),
		Attributes: event.EventAttributes{
			State:     string(s.GetState()),
			PrevState: string(prevState),
		},
	}
}

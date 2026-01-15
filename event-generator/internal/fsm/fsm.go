package fsm

import (
	"event-generator/internal/event"
	"fmt"
	"math/rand"
)

// =======================================================
// Session Interface (FSM이 요구하는 최소 계약)
// =======================================================
// fsm.go (인터페이스)
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

	// 브라우징 컨텍스트
	SetEventPage(string)  // eventPage 저장 메서드
	GetEventPage() string // eventPage 불러오는 메서드
	SetPageType(string)
	GetPageType() string
	SetBrowsingCountryCategory(string)
	GetBrowsingCountryCategory() string
	SetBrowsingProductCategory(string)
	GetBrowsingProductCategory() string
	ResetBrowsingContext()

	// 상품 세부 선택 시 랜덤 선택 후 최소 정보
	SetLastPicked(productID, category, countryr string)
	GetLastProductID() string
	GetLastCategory() string
	GetLastCountry() string
	GetLastPicked() (productID, category, country string)

	SetLastQuantity(qty int)
	GetLastQuantity() int

	// Paging
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
	rnd *rand.Rand
}

func NewSimpleFSM(r *rand.Rand) *SimpleFSM {
	return &SimpleFSM{rnd: r}
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
	tr := chooseTransition(f.rnd, transitions)
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
		s.SetSearchKeyword(FakeKeyword(f.rnd))
		s.SetPageIndex(1)
	}

	// 8. 이벤트 생성
	return &event.Event{
		EventID:   fmt.Sprintf("evt-%d-%09d", now, rand.Int63n(1_000_000_000)),
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

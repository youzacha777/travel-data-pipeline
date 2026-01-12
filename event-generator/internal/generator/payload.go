package generator

import (
	"event-generator/internal/fsm"
	"event-generator/internal/user"
	"math/rand"
	"time"
)

type PayloadGenerator struct {
	rnd *rand.Rand
}

func NewPayloadGenerator() *PayloadGenerator {
	return &PayloadGenerator{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate : 이벤트 타입과 세션 상태에 따라 payload 생성
func (g *PayloadGenerator) Generate(eventType string, session *user.Session) map[string]any {
	var eventPayload map[string]any
	currState := session.GetState()
	prevState := session.GetPrevState()

	switch eventType {

	// 1. 검색 제출 (최초 검색어 입력 시)
	case string(fsm.EventSearchSubmitted):
		eventPayload = g.genSearch(session, eventType)

	// 2. 페이지 조회 (가장 빈번한 이벤트)
	case string(fsm.EventPageViewed):
		// 현재 상태(위치)가 검색 다음 페이지라면 genNextPage 호출
		if currState == fsm.StateNextPage {
			eventPayload = g.genNextPage(session, eventType)
		} else {
			// 그 외 일반적인 페이지 조회는 genBrowsing 호출
			eventPayload = g.genBrowsing(session, eventType)
		}

	// 3. 상품 클릭 (검색 결과 리스트에서 클릭했는지 판별)
	case string(fsm.EventProductClicked):
		// 이전 상태가 검색이었다면 상품 정보를 세션에 저장해야 하므로 genSearch 호출
		if prevState == fsm.StateSearch || prevState == fsm.StateNextPage {
			eventPayload = g.genSearch(session, eventType)
		} else {
			// 홈 화면이나 다른 곳에서 클릭했다면 일반 클릭 로직 수행
			eventPayload = g.genClick(session, eventType)
		}

	// 4. 이벤트 페이지 또는 카테고리 클릭
	case string(fsm.EventPageClicked), string(fsm.EvenCategoryClicked):
		eventPayload = g.genBrowsing(session, eventType)

	// 5. 장바구니 담기 및 뒤로가기
	case string(fsm.EventAddToCart), string(fsm.EventBack):
		eventPayload = g.genClick(session, eventType)

	// 6. 결제 완료 및 최종 종료
	case string(fsm.EventPurchased), string(fsm.EventExit):
		eventPayload = g.genPurchase(session, eventType)

	default:
		// 정의되지 않은 이벤트의 경우 기본 페이로드 생성
		eventPayload = map[string]any{
			"action": "general_step",
		}
	}

	// 공통 데이터 주입 (모든 로그에 필수로 포함되어야 할 정보)
	if eventPayload != nil {
		eventPayload["session_id"] = session.GetID()
		eventPayload["user_id"] = session.GetUserID()
		eventPayload["generated_at"] = session.GetLastEventTs()
		eventPayload["current_state"] = string(currState)
	}

	return eventPayload
}

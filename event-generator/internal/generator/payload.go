package generator

import (
	"event-generator/internal/fsm"
	"event-generator/internal/user"
)

type PayloadGenerator struct {
	// rand/v2 전역 함수를 사용하므로 필드가 필요 없습니다.
}

func NewPayloadGenerator() *PayloadGenerator {
	return &PayloadGenerator{}
}

// Generate : 이벤트 타입과 세션 상태에 따라 payload 생성
func (g *PayloadGenerator) Generate(eventType string, session *user.Session) map[string]any {
	var eventPayload map[string]any
	currState := session.GetState()
	prevState := session.GetPrevState()

	switch eventType {

	// 1. 검색 제출
	case string(fsm.EventSearchSubmitted):
		eventPayload = g.genSearch(session, eventType)

	// 2. 페이지 조회
	case string(fsm.EventPageViewed):
		if currState == fsm.StateNextPage {
			eventPayload = g.genNextPage(session, eventType)
		} else {
			eventPayload = g.genBrowsing(session, eventType)
		}

	// 3. 상품 클릭
	case string(fsm.EventProductClicked):
		if prevState == fsm.StateSearch || prevState == fsm.StateNextPage {
			eventPayload = g.genSearch(session, eventType)
		} else {
			eventPayload = g.genClick(session, eventType)
		}

	// 4. 이벤트 페이지 또는 카테고리 클릭
	case string(fsm.EventPageClicked), string(fsm.EventCategoryClicked):
		eventPayload = g.genBrowsing(session, eventType)

	// 5. 장바구니 담기 및 뒤로가기
	case string(fsm.EventAddToCart):
		eventPayload = g.genClick(session, eventType)

	case string(fsm.EventBack):
		// [수정] 현재 상태가 EventBrowsing이면 전용 로직 호출
		if currState == fsm.StateEventBrowsing {
			eventPayload = g.genEventBrowsing(session, eventType)
		} else {
			eventPayload = g.genClick(session, eventType)
		}

	// 6. 결제 완료 및 최종 종료
	case string(fsm.EventPurchased):
		eventPayload = g.genPurchase(session, eventType)

	case string(fsm.EventExit):
		// EventBrowsing이면 전용 로직 호출
		if currState == fsm.StateEventBrowsing {
			eventPayload = g.genEventBrowsing(session, eventType)
		} else if currState == fsm.StatePurchase {
			eventPayload = g.genPurchase(session, eventType)
		} else {
			eventPayload = g.genBrowsing(session, eventType)
		}

	default:
		eventPayload = map[string]any{
			"action": "general_step",
		}
	}

	// 공통 데이터 주입
	if eventPayload != nil {
		eventPayload["session_id"] = session.GetID()
		eventPayload["user_id"] = session.GetUserID()
		eventPayload["generated_at"] = session.GetLastEventTs()
		eventPayload["current_state"] = string(currState)
	}

	return eventPayload
}

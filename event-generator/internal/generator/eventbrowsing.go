package generator

import (
	"event-generator/internal/fsm"
	"math/rand/v2" // v1 대신 v2를 사용합니다.
)

// genEventBrowsing
func (g *PayloadGenerator) genEventBrowsing(session fsm.Session, eventType string) map[string]any {
	payload := map[string]any{}

	// 현재 사용자가 어떤 이벤트 페이지에 머물고 있는지 세션에서 가져옴
	currentPage := session.GetEventPage()

	switch eventType {

	case string(fsm.EventBack):
		// 이전 페이지로 돌아감 (보통 홈 화면/Browsing 상태로 복귀)
		payload["from_page"] = currentPage
		payload["to_page"] = "home"
		payload["action"] = "back_button_click"

		// [수정] g.rnd.Intn -> rand.IntN (v2 전역 함수 사용하여 Thread-safe 확보)
		// 이벤트 페이지에서 얼마나 머물다 돌아갔는지 기록
		payload["stay_sec"] = rand.IntN(60) + 2

	case string(fsm.EventExit):
		// 앱 종료 혹은 이탈
		payload["last_viewed_page"] = currentPage
		payload["exit_reason"] = "user_left"

		// [수정] g.rnd.Intn -> rand.IntN
		// 이탈 전 최종 체류 시간
		payload["total_event_stay_sec"] = rand.IntN(50) + 10

	}

	return payload
}

package generator

import (
	"event-generator/internal/fsm"
	"math/rand/v2"
)

// GenerateClickPayload
// Click 상태 진입 시 payload 생성
func (g *PayloadGenerator) genClick(session fsm.Session, eventType string) map[string]any {
	payload := map[string]any{}

	// 1. 공통 페이로드: 세션에서 마지막으로 픽한 상품 정보 가져오기
	lastProductID, lastCategory, lastCountry := session.GetLastPicked()

	// 이 정보는 클릭 이후 어떤 이벤트가 발생하든 상세 페이지 로그에는 기본으로 포함
	payload["product_id"] = lastProductID
	payload["category"] = lastCategory
	payload["country"] = lastCountry

	switch eventType {
	case string(fsm.EventAddToCart), string(fsm.EventPurchased):
		// 수량 랜덤 생성 (1~5개)
		quantity := rand.IntN(5) + 1

		// 세션에 수량 저장 (추후 결제 단계 등에서 활용)
		session.SetLastQuantity(quantity)

		// 페이로드에도 현재 행동의 수량 포함
		payload["quantity"] = quantity

	case string(fsm.EventBack):
		// 이전 상태로 돌아가므로 payload 그대로 유지
		payload["stay_sec"] = rand.IntN(30) + 5

	case string(fsm.EventExit):
		payload["exit_reason"] = "user_left"
	}

	return payload
}

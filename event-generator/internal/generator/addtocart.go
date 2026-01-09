package generator

import (
	"event-generator/internal/fsm"
)

// genAddToCart generates payload for shopping cart state events.
func (g *PayloadGenerator) genAddToCart(session fsm.Session, eventType string) map[string]any {
	payload := map[string]any{}

	// 1. 공통 페이로드: 세션에서 상품 정보 및 수량 불러오기
	// 상세 페이지(Click)에서 저장했던 LastPicked와 LastQuantity를 활용합니다.
	lastProductID, lastCategory, lastCountry := session.GetLastPicked()
	lastQuantity := session.GetLastQuantity()

	if lastCategory == "" {
		lastCategory = "unknown"
	}

	payload["product_id"] = lastProductID
	payload["product_category"] = lastCategory
	payload["country"] = lastCountry
	payload["quantity"] = lastQuantity // Click 단계에서 생성된 수량

	// 기본 체류 시간 (장바구니 확인 시간)
	payload["stay_sec"] = g.rnd.Intn(20) + 5

	// 2. 이벤트별 분기 처리
	switch eventType {
	case string(fsm.EventPurchased):
		// 장바구니에서 바로 구매로 넘어가는 경우
		// 분석 시 '상세페이지 직구매'와 '장바구니 경유 구매'를 구분하기 위한 필드
		payload["purchase_source"] = "cart_checkout"

	case string(fsm.EventBack):
		// 장바구니에서 다시 상품 리스트나 상세로 돌아감
		payload["action"] = "back_to_previous"

	case string(fsm.EventExit):
		// 장바구니에 담아만 두고 앱을 종료 (Cart Abandonment 분석 대상)
		payload["exit_reason"] = "user_left_with_items"
	}

	return payload
}

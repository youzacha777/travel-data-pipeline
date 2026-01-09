package generator

import (
	"event-generator/internal/fsm"
)

// genPurchase generates payload for Purchase state events.
// eventType을 인자에 추가하여 switch 문이 작동하도록 수정했습니다.
func (g *PayloadGenerator) genPurchase(session fsm.Session, eventType string) map[string]any {
	payload := map[string]any{}

	// 1. 세션에서 상세 페이지(Click) 단계 때 저장했던 정보들 가져오기
	lastProductID, lastCategory, lastCountry := session.GetLastPicked()
	lastQuantity := session.GetLastQuantity()

	// 2. 공통 상품 정보 구성 (어떤 이벤트든 구매 세션의 상품 정보는 포함)
	payload["product_id"] = lastProductID
	payload["product_category"] = lastCategory
	payload["country"] = lastCountry
	payload["quantity"] = lastQuantity

	// 3. 결제 특화 정보 (기본적으로 결제 성공 페이지 진입 시 생성)
	paymentMethods := []string{"card", "kakao_pay", "naver_pay", "apple_pay", "google_pay"}
	payload["payment_method"] = paymentMethods[g.rnd.Intn(len(paymentMethods))]

	// 4. 기본 체류 시간
	payload["stay_sec"] = g.rnd.Intn(40) + 20

	// 5. 이벤트 타입별 추가 처리
	switch eventType {
	case string(fsm.EventExit):
		// [결제 완료 후 앱 종료 시점]
		payload["action"] = "order_complete_exit"
		payload["exit_reason"] = "user_closed_after_purchase"
		// 종료 시점에는 체류 시간을 조금 짧게 조정 (선택 사항)
		payload["stay_sec"] = g.rnd.Intn(10) + 2
	}

	return payload
}

package generator

import (
	"event-generator/internal/fsm"
	"math/rand/v2" // v1 대신 v2를 사용합니다.
)

// genPurchase
func (g *PayloadGenerator) genPurchase(session fsm.Session, eventType string) map[string]any {
	payload := map[string]any{}

	// 1. 세션에서 상세 페이지(Click) 단계 때 저장했던 정보들 가져오기
	lastProductID, lastCategory, lastCountry := session.GetLastPicked()
	lastQuantity := session.GetLastQuantity()

	// 2. 공통 상품 정보 구성
	payload["product_id"] = lastProductID
	payload["product_category"] = lastCategory
	payload["country"] = lastCountry
	payload["quantity"] = lastQuantity

	// 3. 결제 특화 정보
	paymentMethods := []string{"card", "kakao_pay", "naver_pay", "apple_pay", "google_pay"}
	// [수정] g.rnd.Intn -> rand.IntN (v2 전역 함수 사용)
	payload["payment_method"] = paymentMethods[rand.IntN(len(paymentMethods))]

	// 4. 체류 시간
	// [수정] g.rnd.Intn -> rand.IntN
	payload["stay_sec"] = rand.IntN(40) + 20

	// 5. 이벤트 타입별 추가 처리
	switch eventType {
	case string(fsm.EventExit):
		payload["action"] = "order_complete_exit"
		payload["exit_reason"] = "user_closed_after_purchase"
		// [수정] g.rnd.Intn -> rand.IntN
		payload["stay_sec"] = rand.IntN(10) + 2
	}

	return payload
}

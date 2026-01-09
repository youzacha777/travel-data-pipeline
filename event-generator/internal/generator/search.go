package generator

import (
	"event-generator/internal/fsm"
	"log"
)

// genSearch generates payload for Search state events.
func (g *PayloadGenerator) genSearch(session fsm.Session, eventType string) map[string]any {
	payload := map[string]any{
		"query": session.GetSearchKeyword(),
	}

	switch eventType {

	case string(fsm.EventPageViewed):
		// 검색 결과 페이지 탐색
		payload["stay_sec"] = g.rnd.Intn(40) + 1

	case string(fsm.EventBack):
		// 홈으로 뒤로가기
		payload["exit_reason"] = "back_to_home"
		payload["stay_sec"] = g.rnd.Intn(5) + 1

	case string(fsm.EventExit):
		// 검색 이탈
		payload["exit_reason"] = "search_exit"
		payload["stay_sec"] = g.rnd.Intn(5) + 1

	case string(fsm.EventProductClicked):
		keyword := session.GetSearchKeyword()
		// 1. 키워드로 상품 찾기
		product, _ := DistinguishAndGetProduct(keyword)

		if product != nil {
			// 2. 세션에 저장 (이게 되어야 다음 Purchase에서 꺼낼 수 있음)
			session.SetLastPicked(product.ProductID, product.Category, product.Country)

			// 3. 페이로드 구성
			payload["product_id"] = product.ProductID
			payload["category"] = product.Category
			payload["country"] = product.Country
		} else {
			// [중요] 상품을 못 찾았을 때 로그를 찍어봐야 합니다.
			// 데이터 맵에 해당 키워드에 맞는 상품이 없는 경우입니다.
			log.Printf("[WARN] No product matched for keyword: '%s'", keyword)

			// 방어 로직: 상품을 못 찾으면 임의의 상품이라도 세팅할지 결정해야 합니다.
			// 현재는 비어있는 상태로 넘어가므로 이후 단계에서 0이나 빈값이 찍히게 됩니다.
		}

	}

	return payload
}

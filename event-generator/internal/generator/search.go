package generator

import (
	"event-generator/internal/fsm"
	"log"
	"math/rand/v2" // v1 대신 v2를 사용합니다.
)

// genSearch
func (g *PayloadGenerator) genSearch(session fsm.Session, eventType string) map[string]any {
	payload := map[string]any{
		"query": session.GetSearchKeyword(),
	}

	switch eventType {

	case string(fsm.EventPageViewed):
		// 검색 결과 페이지 탐색
		// [수정] g.rnd.Intn -> rand.IntN (v2 전역 함수 사용)
		payload["stay_sec"] = rand.IntN(40) + 1

	case string(fsm.EventBack):
		// 홈으로 뒤로가기
		payload["exit_reason"] = "back_to_home"
		// [수정] g.rnd.Intn -> rand.IntN
		payload["stay_sec"] = rand.IntN(5) + 1

	case string(fsm.EventExit):
		// 검색 이탈
		payload["exit_reason"] = "search_exit"
		// [수정] g.rnd.Intn -> rand.IntN
		payload["stay_sec"] = rand.IntN(5) + 1

	case string(fsm.EventProductClicked):
		keyword := session.GetSearchKeyword()
		// 키워드로 상품 구분 및 획득
		product, _ := DistinguishAndGetProduct(keyword)

		if product != nil {
			// 세션에 저장
			session.SetLastPicked(product.ProductID, product.Category, product.Country)

			// 3. 페이로드 구성
			payload["product_id"] = product.ProductID
			payload["category"] = product.Category
			payload["country"] = product.Country
		} else {
			// 상품을 못 찾았을 때 로그
			log.Printf("[WARN] No product matched for keyword: '%s'", keyword)
		}
	}

	return payload
}

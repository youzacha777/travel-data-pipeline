package generator

import (
	"event-generator/internal/fsm"
	"log"
)

// GenerateNextPagePayload
// NextPage 상태 진입 시 payload 생성
func (g *PayloadGenerator) genNextPage(session fsm.Session, eventType string) map[string]any {
	payload := map[string]any{}

	// 기본 검색 컨텍스트 유지
	payload["query"] = session.GetSearchKeyword()

	// 페이지 인덱스 증가 (NextPage 상태 진입 시 항상 증가)
	session.IncrementPageIndex()
	payload["page_index"] = session.GetPageIndex()

	// 체류 시간
	payload["stay_sec"] = g.rnd.Intn(40) + 10

	switch eventType {
	case string(fsm.EventPageViewed):
		// [계속 스크롤/다음 페이지 조회]
		// 페이지 인덱스 증가
		session.IncrementPageIndex()

		// 증가된 페이지 번호를 페이로드에 할당
		payload["page_index"] = session.GetPageIndex()
		payload["action"] = "scroll_next_page"

	case string(fsm.EventProductClicked):
		// 상품명 기반 검색
		keyword := session.GetSearchKeyword()

		// 1. 상품/국가/카테고리 판별 함수 사용
		product, searchType := DistinguishAndGetProduct(keyword)

		// 2. 방어 로직: 검색 결과가 아예 없는 경우
		if product == nil {
			log.Printf("[WARN] No product found for keyword: %s (Type: %s)\n", keyword, searchType)
			return nil
		}

		// 3. 세션 업데이트 (제공해주신 SetLastPicked 활용)
		session.SetLastPicked(product.ProductID, product.Category, product.Country)

		// 4. 페이로드 구성 (추천 분류값 포함)
		payload["product_id"] = product.ProductID
		payload["product_name"] = product.ProductName
		payload["category"] = product.Category
		payload["country"] = product.Country
		payload["search_type"] = searchType // 검색어가 뭐였는지 분석용 데이터 추가

	case string(fsm.EventBack):
		// 이전 페이지로 돌아감
		payload["action"] = "back_button_click"

		// 이벤트 페이지에서 얼마나 머물다 돌아갔는지 기록
		payload["stay_sec"] = g.rnd.Intn(60) + 2

	case string(fsm.EventExit):
		// 앱 종료 혹은 이탈
		payload["exit_reason"] = "user_left"

		// 이탈 전 최종 체류 시간
		payload["total_event_stay_sec"] = g.rnd.Intn(50) + 10

	}

	return payload
}

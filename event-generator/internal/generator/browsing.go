package generator

import (
	"event-generator/internal/fsm"
)

// genBrowsing generates payload for browsing-related events based on session state.
func (g *PayloadGenerator) genBrowsing(session fsm.Session, eventType string) map[string]any {
	payload := map[string]any{}

	switch eventType {

	case string(fsm.EventSearchSubmitted):
		payload["query"] = session.GetSearchKeyword()
		payload["stay_sec"] = g.rnd.Intn(10) + 1

	case string(fsm.EventPageViewed):
		session.SetPageType("first_page")
		payload["page_type"] = "first_page"
		payload["stay_sec"] = g.rnd.Intn(180) + 5

	case string(fsm.EventPageClicked):
		pageTypes := []string{"special_event_category", "recommend_category"}
		pageType := pageTypes[g.rnd.Intn(len(pageTypes))]
		session.SetPageType(pageType)
		payload["page_type"] = pageType

		switch pageType {
		case "special_event_category":
			eventPages := []string{"flight_promotion", "referral_promotion", "continent_promotion", "season_promotion"}
			eventPage := eventPages[g.rnd.Intn(len(eventPages))]
			session.SetEventPage(eventPage)
			payload["special_event_category"] = eventPage
		case "recommend_category":
			session.SetEventPage("recommend_category")
			payload["recommend_category"] = "recommend_list_to_friends"
		}
		payload["stay_sec"] = g.rnd.Intn(180) + 5

	case string(fsm.EventProductClicked):
		// s -> session으로 수정, r -> g.rnd로 수정
		product := pickTopProduct(g.rnd)
		session.SetLastPicked(product.ProductID, product.Category, product.Country)

		payload["product_id"] = product.ProductID
		payload["product_name"] = product.ProductName
		payload["category"] = product.Category
		payload["country"] = product.Country
		payload["recommend_category"] = "recommend_list_to_friends" // 아까 요청하신 필드 추가
		payload["stay_sec"] = g.rnd.Intn(180) + 5

	case string(fsm.EvenCategoryClicked): // Even -> Event 오타 주의 (상수 확인 필요)
		pageTypes := []string{"country_category", "product_category"}
		pageType := pageTypes[g.rnd.Intn(len(pageTypes))]
		session.SetPageType(pageType)
		payload["page_type"] = pageType

		switch pageType {
		case "country_category":
			countries := []string{
				CountryHongKong, CountryTaiwan, CountryMacau, CountrySingapore,
				CountryMalaysia, CountryThailand, CountryUAE, CountryUSA,
			}
			selectedCountry := countries[g.rnd.Intn(len(countries))]
			if product, ok := GetRandomProductByCountry(selectedCountry); ok {
				session.SetLastPicked(product.ProductID, product.Category, product.Country)
				payload["selected_country"] = selectedCountry
				payload["product_id"] = product.ProductID
				payload["product_name"] = product.ProductName
				payload["category"] = product.Category
				payload["recommend_category"] = "country_navigation_list"
			}

		case "product_category":
			categories := []string{
				CategoryAttraction, CategoryTransport, CategoryMuseum,
				CategoryFood, CategoryTour, CategoryShow,
				CategoryExhibition, CategoryEtc,
			}
			selectedCategory := categories[g.rnd.Intn(len(categories))]
			if product, ok := GetRandomProductByCategory(selectedCategory); ok {
				session.SetLastPicked(product.ProductID, product.Category, product.Country)
				payload["selected_category"] = selectedCategory
				payload["product_id"] = product.ProductID
				payload["product_name"] = product.ProductName
				payload["country"] = product.Country
				payload["recommend_category"] = "category_navigation_list"
			}
		}
		payload["stay_sec"] = g.rnd.Intn(180) + 5

	case string(fsm.EventExit):
		payload["exit_reason"] = "user_left"
	}

	return payload
}

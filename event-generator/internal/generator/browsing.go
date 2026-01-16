package generator

import (
	"event-generator/internal/fsm"
	"math/rand/v2" // v1 대신 v2를 사용합니다.
)

// genBrowsing
func (g *PayloadGenerator) genBrowsing(session fsm.Session, eventType string) map[string]any {
	payload := map[string]any{}

	switch eventType {

	case string(fsm.EventSearchSubmitted):
		payload["query"] = session.GetSearchKeyword()
		// [수정] rand.IntN 사용
		payload["stay_sec"] = rand.IntN(10) + 1

	case string(fsm.EventPageViewed):
		session.SetPageType("first_page")
		payload["page_type"] = "first_page"
		// [수정] rand.IntN 사용
		payload["stay_sec"] = rand.IntN(180) + 5

	case string(fsm.EventPageClicked):
		pageTypes := []string{"special_event_category", "recommend_category"}
		// [수정] rand.IntN 사용
		pageType := pageTypes[rand.IntN(len(pageTypes))]
		session.SetPageType(pageType)
		payload["page_type"] = pageType

		switch pageType {
		case "special_event_category":
			eventPages := []string{"flight_promotion", "referral_promotion", "continent_promotion", "season_promotion"}
			eventPage := eventPages[rand.IntN(len(eventPages))]
			session.SetEventPage(eventPage)
			payload["special_event_category"] = eventPage
		case "recommend_category":
			session.SetEventPage("recommend_category")
			payload["recommend_category"] = "recommend_list_to_friends"
		}
		payload["stay_sec"] = rand.IntN(180) + 5

	case string(fsm.EventProductClicked):
		// [핵심 수정] pickTopProduct는 이제 인자를 받지 않습니다.
		// g.rnd를 넘기지 않고 바로 호출합니다.
		product := pickTopProduct()
		session.SetLastPicked(product.ProductID, product.Category, product.Country)

		payload["product_id"] = product.ProductID
		payload["product_name"] = product.ProductName
		payload["category"] = product.Category
		payload["country"] = product.Country
		payload["stay_sec"] = rand.IntN(180) + 5

	case string(fsm.EventCategoryClicked):
		pageTypes := []string{"country_category", "product_category"}
		// [수정] rand.IntN 사용
		pageType := pageTypes[rand.IntN(len(pageTypes))]
		session.SetPageType(pageType)
		payload["page_type"] = pageType

		switch pageType {
		case "country_category":
			countries := []string{
				CountryHongKong, CountryTaiwan, CountryMacau, CountrySingapore,
				CountryMalaysia, CountryThailand, CountryUAE, CountryUSA,
			}
			selectedCountry := countries[rand.IntN(len(countries))]
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
			selectedCategory := categories[rand.IntN(len(categories))]
			if product, ok := GetRandomProductByCategory(selectedCategory); ok {
				session.SetLastPicked(product.ProductID, product.Category, product.Country)
				payload["selected_category"] = selectedCategory
				payload["product_id"] = product.ProductID
				payload["product_name"] = product.ProductName
				payload["country"] = product.Country
				payload["recommend_category"] = "category_navigation_list"
			}
		}
		payload["stay_sec"] = rand.IntN(180) + 5

	case string(fsm.EventExit):
		payload["exit_reason"] = "user_left"
	}

	return payload
}

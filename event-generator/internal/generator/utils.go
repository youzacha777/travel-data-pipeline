package generator

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// 검색어 세팅용 함수
func (g *PayloadGenerator) generateKeyword() string {
	keywords := []string{
		// 국가 검색어
		"홍콩", "대만", "마카오", "싱가포르", "말레이시아", "태국", "UAE", "미국",

		// 콤보 상품
		"콤보 상품",

		// 🇭🇰 홍콩
		"홍콩 디즈니", "코타이젯", "홍콩 터보젯", "피크트램", "옹핑 케이블카",

		// 🇹🇼 대만
		"국립박물관", "딘 타이 펑", "타이페이 101", "Easy 심카드",

		// 🇲🇴 마카오
		"마카오 오픈 탑 버스", "마카오 해리 포터", "마카오 터보젯", "타워 360", "마카오 전망대",

		// 🇸🇬 싱가포르
		"가든스 바이 더 베이", "유니버셜 스튜디오 싱가포르", "윙스 오브 타임",
		"싱가포르 플라이어", "리버크루즈", "나이트 사파리",

		// 🇲🇾 말레이시아
		"레고랜드", "슈퍼파크 말레이시아", "썬웨이 라군",

		// 🇹🇭 태국
		"5G 심카드", "진리의 성전", "마하나콘 전망대", "푸켓 아쿠아리움",

		// 🇦🇪 UAE
		"카스르 알 와탄", "페라리 월드 아부다비", "루브르 아부다비",
		"부르즈 할리파", "더 뷰 앳 더 팜", "글로벌 빌리지 두바이",

		// 🇺🇸 미국
		"미국 자연사 박물관", "캘리포니아 디즈니",
		"LA 빅 버스", "MoMA", "탑 오브 더 락",
	}

	return keywords[g.rnd.Intn(len(keywords))]
}

// 국가 카테고리 키워드별 검색어 생성용 함수

var countryProductPool = map[string][]string{
	"홍콩":    {"홍콩 디즈니", "코타이젯", "홍콩 터보젯", "피크트램", "옹핑 케이블카"},
	"대만":    {"대만 국립박물관", "딘 타이 펑", "타이페이 101", "Easy 심카드"},
	"마카오":   {"마카오 오픈 탑 버스", "마카오 해리 포터", "마카오 터보젯", "타워 360", "마카오 전망대"},
	"싱가포르":  {"가든스 바이 더 베이", "유니버셜 스튜디오 싱가포르", "윙스 오브 타임", "싱가포르 플라이어", "리버크루즈", "나이트 사파리"},
	"말레이시아": {"레고랜드", "슈퍼파크 말레이시아", "썬웨이 라군"},
	"태국":    {"5G 심카드", "진리의 성전", "마하나콘 전망대", "푸켓 아쿠아리움"},
	"UAE":   {"카스르 알 와탄", "페라리 월드 아부다비", "루브르 아부다비", "부르즈 할리파", "더 뷰 앳 더 팜", "글로벌 빌리지 두바이"},
	"미국":    {"미국 자연사 박물관", "캘리포니아 디즈니", "LA 빅 버스", "MoMA", "탑 오브 더 락"},
}

func (g *PayloadGenerator) pickTopCountryProducts(country string, n int) []string {
	pool := countryProductPool[country]
	if len(pool) == 0 {
		return nil
	}

	if n >= len(pool) {
		return append([]string{}, pool...)
	}

	return append([]string{}, pool[:n]...)
}

// EventProductClicked 이벤트 발생 시 세부 상품 선택 함수

type PickedProduct struct {
	Product Product
	Vendor  string
}

// Top5 노출 상품 랜덤 선택 함수

// 홈 상단 노출 대상 상품 ID
var homeExposureProductIDs = []string{
	"P001", // 홍콩 디즈니
	"P016", // 유니버셜 스튜디오 싱가포르
	"P022", // 슈퍼파크 말레이시아
	"P029", // 페라리 월드 아부다비
	"P035", // 캘리포니아 디즈니
}

func pickTopExposureProduct(r *rand.Rand) Product {
	pid := homeExposureProductIDs[r.Intn(len(homeExposureProductIDs))]

	for _, p := range products {
		if p.ProductID == pid {
			return p
		}
	}

	// 논리적으로 여기 오면 안 됨
	panic("home exposure product not found: " + pid)
}

// search에서 사용되는 함수

// 국가명에 맞는 상품을 랜덤으로 선택하는 함수
func pickProductByCountry(country string) (Product, error) {
	var candidates []Product

	// 해당 국가에 맞는 상품들을 필터링
	for _, p := range products {
		if p.Country == country {
			candidates = append(candidates, p)
		}
	}

	// 국가에 맞는 상품이 없으면 오류 반환
	if len(candidates) == 0 {
		log.Println("No product found for the given country:", country) // 여기서 로그 출력
		return Product{}, fmt.Errorf("No products found for country: %s", country)
	}

	// 랜덤으로 상품 선택
	rand.Seed(time.Now().UnixNano())
	return candidates[rand.Intn(len(candidates))], nil
}

// 상품 카테고리에 맞는 상품을 랜덤으로 선택하는 함수
func pickProductByCategory(category string) (Product, error) {
	var candidates []Product

	// 해당 카테고리에 맞는 상품들을 필터링
	for _, p := range products {
		if p.Category == category {
			candidates = append(candidates, p)
		}
	}

	// 카테고리에 맞는 상품이 없으면 오류 반환
	if len(candidates) == 0 {
		log.Println("No product found for the given category:", category) // 여기서 로그 출력
		return Product{}, fmt.Errorf("No products found for category: %s", category)
	}

	// 랜덤으로 상품 선택
	return candidates[rand.Intn(len(candidates))], nil
}

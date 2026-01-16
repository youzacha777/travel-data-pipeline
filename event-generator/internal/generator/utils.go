package generator

import (
	"fmt"
	"log"
	"math/rand/v2" // v1 대신 v2를 사용합니다. (Thread-safe 최적화)
)

// 검색어 세팅용 함수
// [수정] 인자에서 g.rnd 의존성 제거
func (g *PayloadGenerator) generateKeyword() string {
	keywords := []string{
		"홍콩", "대만", "마카오", "싱가포르", "말레이시아", "태국", "UAE", "미국",
		"콤보 상품",
		"홍콩 디즈니", "코타이젯", "홍콩 터보젯", "피크트램", "옹핑 케이블카",
		"국립박물관", "딘 타이 펑", "타이페이 101", "Easy 심카드",
		"마카오 오픈 탑 버스", "마카오 해리 포터", "마카오 터보젯", "타워 360", "마카오 전망대",
		"가든스 바이 더 베이", "유니버셜 스튜디오 싱가포르", "윙스 오브 타임",
		"싱가포르 플라이어", "리버크루즈", "나이트 사파리",
		"레고랜드", "슈퍼파크 말레이시아", "썬웨이 라군",
		"5G 심카드", "진리의 성전", "마하나콘 전망대", "푸켓 아쿠아리움",
		"카스르 알 와탄", "페라리 월드 아부다비", "루브르 아부다비",
		"부르즈 할리파", "더 뷰 앳 더 팜", "글로벌 빌리지 두바이",
		"미국 자연사 박물관", "캘리포니아 디즈니",
		"LA 빅 버스", "MoMA", "탑 오브 더 락",
	}

	// [수정] rand.Intn -> rand.IntN (v2 전역 함수)
	return keywords[rand.IntN(len(keywords))]
}

// 국가별 상품 풀 (기존 동일)
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

type PickedProduct struct {
	Product Product
	Vendor  string
}

// 홈 상단 노출 대상 상품 ID
var homeExposureProductIDs = []string{
	"P001", // 홍콩 디즈니
	"P016", // 유니버셜 스튜디오 싱가포르
	"P022", // 슈퍼파크 말레이시아
	"P029", // 페라리 월드 아부다비
	"P035", // 캘리포니아 디즈니
}

// [수정] 인자에서 r *rand.Rand 제거 및 v2 적용
func pickTopExposureProduct() Product {
	// [수정] rand.IntN 사용
	pid := homeExposureProductIDs[rand.IntN(len(homeExposureProductIDs))]

	for _, p := range products {
		if p.ProductID == pid {
			return p
		}
	}

	panic("home exposure product not found: " + pid)
}

// 국가명에 맞는 상품을 랜덤으로 선택하는 함수
func pickProductByCountry(country string) (Product, error) {
	var candidates []Product

	for _, p := range products {
		if p.Country == country {
			candidates = append(candidates, p)
		}
	}

	if len(candidates) == 0 {
		log.Println("No product found for the given country:", country)
		return Product{}, fmt.Errorf("No products found for country: %s", country)
	}

	// [수정] rand.Seed 제거 (v2에서는 필요 없으며 성능 저하 원인임)
	// [수정] rand.IntN 사용
	return candidates[rand.IntN(len(candidates))], nil
}

// 상품 카테고리에 맞는 상품을 랜덤으로 선택하는 함수
func pickProductByCategory(category string) (Product, error) {
	var candidates []Product

	for _, p := range products {
		if p.Category == category {
			candidates = append(candidates, p)
		}
	}

	if len(candidates) == 0 {
		log.Println("No product found for the given category:", category)
		return Product{}, fmt.Errorf("No products found for category: %s", category)
	}

	// [수정] rand.IntN 사용
	return candidates[rand.IntN(len(candidates))], nil
}

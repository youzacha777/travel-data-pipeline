package generator

import (
	"math/rand"
	"strings"
)

type Product struct {
	ProductID   string
	ProductName string
	Country     string
	Category    string
	Vendors     []string
}

// 카테고리 상수
const (
	CategoryAttraction = "attraction"
	CategoryTransport  = "transport"
	CategoryMuseum     = "museum"
	CategoryFood       = "food"
	CategoryTour       = "tour"
	CategoryShow       = "show"
	CategoryExhibition = "exhibition"
	CategoryEtc        = "etc"
)

// 국가명 상수
const (
	CountryHongKong  = "홍콩"
	CountryTaiwan    = "대만"
	CountryMacau     = "마카오"
	CountrySingapore = "싱가포르"
	CountryMalaysia  = "말레이시아"
	CountryThailand  = "태국"
	CountryUAE       = "UAE"
	CountryUSA       = "미국"
)

// 빠른 검색을 위한 전역 변수
var (
	productMap  map[string]*Product
	countryMap  map[string][]*Product
	categoryMap map[string][]*Product
)

// 홈 상단 노출 대상 상품명 리스트
var homeExposureProductNames = []string{
	"홍콩 디즈니",
	"유니버셜 스튜디오 싱가포르",
	"슈퍼파크 말레이시아",
	"페라리 월드 아부다비",
	"캘리포니아 디즈니",
}

// 데이터 초기화
func init() {
	productMap = make(map[string]*Product)
	countryMap = make(map[string][]*Product)
	categoryMap = make(map[string][]*Product) // <-- 추가

	for i := range products {
		p := &products[i]

		// 1. 이름으로 바로 찾기 (1:1)
		productMap[p.ProductName] = p

		// 2. 국가별 상품 묶음 (1:N)
		countryMap[p.Country] = append(countryMap[p.Country], p)

		// 3. 카테고리별 상품 묶음 (1:N) <-- 추가
		categoryMap[p.Category] = append(categoryMap[p.Category], p)

	}

}

// GetProductByName: 상품명으로 정확히 일치하는 상품 정보 반환
func GetProductByName(name string) (*Product, bool) {
	p, ok := productMap[name]
	return p, ok
}

// GetRandomProductByCountry: 국가명으로 검색하여 해당 국가 상품 중 랜덤 1개 반환
func GetRandomProductByCountry(country string) (*Product, bool) {
	products, ok := countryMap[country]
	if !ok || len(products) == 0 {
		return nil, false
	}

	// 랜덤 인덱스 선택
	randomIndex := rand.Intn(len(products))
	return products[randomIndex], true
}

// GetRandomProductByCategory: 카테고리명으로 검색하여 해당 카테고리 상품 중 랜덤 1개 반환
func GetRandomProductByCategory(category string) (*Product, bool) {
	products, ok := categoryMap[category]
	if !ok || len(products) == 0 {
		return nil, false
	}

	randomIndex := rand.Intn(len(products))
	return products[randomIndex], true
}

// pickTopExposureProduct: 홈 노출 리스트 중 하나를 랜덤하게 뽑아 상세 정보 반환
func pickTopProduct(r *rand.Rand) *Product {
	// 1. 이름 리스트에서 랜덤하게 하나 선택
	name := homeExposureProductNames[r.Intn(len(homeExposureProductNames))]

	// 2. 이미 만들어둔 GetProductByName 함수 활용
	p, ok := GetProductByName(name)
	if !ok {
		// init에서 맵핑이 잘 되었다면 절대 발생하지 않음
		panic("exposure product name not found in map: " + name)
	}

	return p
}

func DistinguishAndGetProduct(query string) (*Product, string) {
	// 1. [상품명 완전 일치] 우선 확인
	// "홍콩 디즈니"라고 쳤을 때 "홍콩" 국가로 빠지지 않고 해당 상품을 바로 주기 위함
	if p, ok := GetProductByName(query); ok {
		return p, "product_match"
	}

	// 2. [국가명 완전 일치] 확인
	// 키워드가 우리가 정의한 국가 상수(예: "태국", "미국")와 정확히 일치하는지 확인
	if _, ok := countryMap[query]; ok {
		p, found := GetRandomProductByCountry(query)
		if found {
			return p, "country_match"
		}
	}

	// 3. [카테고리명 완전 일치] 확인
	// 키워드가 "museum", "attraction" 등 카테고리명인지 확인
	if _, ok := categoryMap[query]; ok {
		p, found := GetRandomProductByCategory(query)
		if found {
			return p, "category_match"
		}
	}

	// 4. [부분 일치] (Optional)
	// 위에서 완전히 일치하는 키워드가 없을 때만 마지막으로 이름에 포함되어 있는지 확인
	for name, p := range productMap {
		if strings.Contains(name, query) {
			return p, "partial_match"
		}
	}

	return nil, "no_match"
}

// 상품 원본 데이터
var products = []Product{
	{ProductID: "P001", ProductName: "홍콩 디즈니", Country: CountryHongKong, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorB", "VendorC"}},
	{ProductID: "P002", ProductName: "코타이젯", Country: CountryHongKong, Category: CategoryTransport, Vendors: []string{"VendorA", "VendorC"}},
	{ProductID: "P003", ProductName: "홍콩 터보젯", Country: CountryHongKong, Category: CategoryTransport, Vendors: []string{"VendorA", "VendorC"}},
	{ProductID: "P004", ProductName: "피크트램", Country: CountryHongKong, Category: CategoryTransport, Vendors: []string{"VendorA", "VendorB"}},
	{ProductID: "P005", ProductName: "옹핑 케이블카", Country: CountryHongKong, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorB", "VendorC"}},
	{ProductID: "P006", ProductName: "대만 국립박물관", Country: CountryTaiwan, Category: CategoryMuseum, Vendors: []string{"VendorA", "VendorC"}},
	{ProductID: "P007", ProductName: "딘 타이 펑", Country: CountryTaiwan, Category: CategoryFood, Vendors: []string{"VendorB", "VendorC"}},
	{ProductID: "P008", ProductName: "타이페이 101", Country: CountryTaiwan, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorB"}},
	{ProductID: "P009", ProductName: "Easy 심카드", Country: CountryTaiwan, Category: CategoryEtc, Vendors: []string{"VendorC"}},
	{ProductID: "P010", ProductName: "마카오 오픈 탑 버스", Country: CountryMacau, Category: CategoryTransport, Vendors: []string{"VendorA", "VendorB"}},
	{ProductID: "P011", ProductName: "마카오 해리 포터", Country: CountryMacau, Category: CategoryExhibition, Vendors: []string{"VendorA"}},
	{ProductID: "P012", ProductName: "마카오 터보젯", Country: CountryMacau, Category: CategoryTransport, Vendors: []string{"VendorA", "VendorC"}},
	{ProductID: "P013", ProductName: "타워 360", Country: CountryMacau, Category: CategoryAttraction, Vendors: []string{"VendorB", "VendorC"}},
	{ProductID: "P014", ProductName: "마카오 전망대", Country: CountryMacau, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorB"}},
	{ProductID: "P015", ProductName: "가든스 바이 더 베이", Country: CountrySingapore, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorB", "VendorC"}},
	{ProductID: "P016", ProductName: "유니버셜 스튜디오 싱가포르", Country: CountrySingapore, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorC"}},
	{ProductID: "P017", ProductName: "윙스 오브 타임", Country: CountrySingapore, Category: CategoryShow, Vendors: []string{"VendorA"}},
	{ProductID: "P018", ProductName: "싱가포르 플라이어", Country: CountrySingapore, Category: CategoryAttraction, Vendors: []string{"VendorB"}},
	{ProductID: "P019", ProductName: "리버크루즈", Country: CountrySingapore, Category: CategoryTour, Vendors: []string{"VendorA", "VendorC"}},
	{ProductID: "P020", ProductName: "나이트 사파리", Country: CountrySingapore, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorB"}},
	{ProductID: "P021", ProductName: "레고랜드", Country: CountryMalaysia, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorC"}},
	{ProductID: "P022", ProductName: "슈퍼파크 말레이시아", Country: CountryMalaysia, Category: CategoryAttraction, Vendors: []string{"VendorB"}},
	{ProductID: "P023", ProductName: "5G 심카드", Country: CountryMalaysia, Category: CategoryEtc, Vendors: []string{"VendorC"}},
	{ProductID: "P024", ProductName: "썬웨이 라군", Country: CountryMalaysia, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorB"}},
	{ProductID: "P025", ProductName: "진리의 성전", Country: CountryThailand, Category: CategoryAttraction, Vendors: []string{"VendorA"}},
	{ProductID: "P026", ProductName: "마하나콘 전망대", Country: CountryThailand, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorB"}},
	{ProductID: "P027", ProductName: "푸켓 아쿠아리움", Country: CountryThailand, Category: CategoryAttraction, Vendors: []string{"VendorC"}},
	{ProductID: "P028", ProductName: "카스르 알 와탄", Country: CountryUAE, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorB"}},
	{ProductID: "P029", ProductName: "페라리 월드 아부다비", Country: CountryUAE, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorC"}},
	{ProductID: "P030", ProductName: "루브르 아부다비", Country: CountryUAE, Category: CategoryMuseum, Vendors: []string{"VendorB", "VendorC"}},
	{ProductID: "P031", ProductName: "부르즈 할리파", Country: CountryUAE, Category: CategoryMuseum, Vendors: []string{"VendorA", "VendorB", "VendorC"}},
	{ProductID: "P032", ProductName: "더 뷰 앳 더 팜", Country: CountryUAE, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorC"}},
	{ProductID: "P033", ProductName: "글로벌 빌리지 두바이", Country: CountryUAE, Category: CategoryAttraction, Vendors: []string{"VendorB"}},
	{ProductID: "P034", ProductName: "미국 자연사 박물관", Country: CountryUSA, Category: CategoryMuseum, Vendors: []string{"VendorA"}},
	{ProductID: "P035", ProductName: "캘리포니아 디즈니", Country: CountryUSA, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorC"}},
	{ProductID: "P036", ProductName: "LA 빅 버스 투어", Country: CountryUSA, Category: CategoryTour, Vendors: []string{"VendorB"}},
	{ProductID: "P037", ProductName: "MoMA 현대 미술관", Country: CountryUSA, Category: CategoryMuseum, Vendors: []string{"VendorC"}},
	{ProductID: "P038", ProductName: "탑 오브 더 락", Country: CountryUSA, Category: CategoryAttraction, Vendors: []string{"VendorA", "VendorB"}},
}

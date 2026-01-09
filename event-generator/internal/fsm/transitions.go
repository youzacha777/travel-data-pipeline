package fsm

type Transition struct {
	Event     EventType
	NextState State
	Weight    float64
}

// Transitions defines the user behavior model.
// - Key: current state
// - Value: possible transitions from that state
var Transitions = map[State][]Transition{

	// =========================================================
	// Level 1: Browsing (초기 탐색)
	// =========================================================
	StateBrowsing: {
		// 검색 진입
		{Event: EventSearchSubmitted, NextState: StateSearch, Weight: 0.3},

		// 노출 상품 클릭
		{Event: EventProductClicked, NextState: StateClick, Weight: 0.2}, // 노출 상품 클릭으로 조건 설정 필요

		// 이벤트 페이지 클릭
		{Event: EventPageClicked, NextState: StateEventBrowsing, Weight: 0.2}, // 이벤트 페이지 클릭시 상세 조건 별 분기는 fsm.step 함수에서 처리

		// 카테고리 페이지 클릭
		{Event: EvenCategoryClicked, NextState: StateClick, Weight: 0.2}, // 이벤트 페이지 클릭시 상세 조건 별 분기는 fsm.step 함수에서 처리

		// 홈/첫 페이지 조회 (머무름)
		{Event: EventPageViewed, NextState: StateBrowsing, Weight: 0.5}, // 상태 변화 없이 머무름

		// 이탈
		{Event: EventExit, NextState: StateExit, Weight: 0.5},
	},

	// =========================================================
	// Level 2: EventBrowsing (이벤트 탐색)
	// =========================================================

	StateEventBrowsing: {

		// 뒤로 → 홈 (back 이벤트는 Step()에서 PrevState로 override됨)
		{Event: EventBack, NextState: "", Weight: 0.7},

		// 이탈
		{Event: EventExit, NextState: StateExit, Weight: 0.3},
	},

	// =========================================================
	// Level 3: Search (의도 형성)
	// =========================================================
	StateSearch: {
		// 검색 결과 상세 진입
		{Event: EventProductClicked, NextState: StateClick, Weight: 0.6},

		// 검색 결과 페이지 탐색(머무르거나 페이지 넘기기)
		{Event: EventPageViewed, NextState: StateNextPage, Weight: 0.1},

		// 뒤로 → 홈 (back 이벤트는 Step()에서 PrevState로 override됨)
		{Event: EventBack, NextState: "", Weight: 0.2},

		// 이탈
		{Event: EventExit, NextState: StateExit, Weight: 0.1},
	},

	// =========================================================
	// Level 2: NextPage (탐색 심화)
	// =========================================================
	StateNextPage: {
		// 상품 상세 진입
		{Event: EventProductClicked, NextState: StateClick, Weight: 0.6},

		// 계속 스크롤
		{Event: EventPageViewed, NextState: StateNextPage, Weight: 0.2},

		// 뒤로 → 검색 (back 이벤트는 Step()에서 PrevState로 override됨)
		{Event: EventBack, NextState: "", Weight: 0.1},

		// 이탈
		{Event: EventExit, NextState: StateExit, Weight: 0.1},
	},

	// =========================================================
	// Level 3: Click (상품 상세)
	// =========================================================
	StateClick: {
		// 장바구니
		{Event: EventAddToCart, NextState: StateAddToCart, Weight: 0.4},

		// 바로 구매
		{Event: EventPurchased, NextState: StatePurchase, Weight: 0.4},

		// 뒤로 → 탐색 (back 이벤트는 Step()에서 PrevState로 override됨)
		{Event: EventBack, NextState: "", Weight: 0.1},

		// 이탈
		{Event: EventExit, NextState: StateExit, Weight: 0.1},
	},

	// =========================================================
	// Level 3: AddToCart (전환 직전)
	// =========================================================
	StateAddToCart: {
		// 구매 확정
		{Event: EventPurchased, NextState: StatePurchase, Weight: 0.7},

		// 뒤로 → 상세
		{Event: EventBack, NextState: "", Weight: 0.2},

		// 이탈
		{Event: EventExit, NextState: StateExit, Weight: 0.1},
	},

	// =========================================================
	// Level 4: Terminal States
	// =========================================================
	StatePurchase: {{Event: EventExit, NextState: StateExit, Weight: 1}},
	StateExit:     {},
}

// fsm/types.go

package fsm

type State string
type EventType string

const (
	// States (위치)
	StateBrowsing      State = "browsing"
	StateEventBrowsing State = "eventbrowsing"
	StateSearch        State = "search"
	StateNextPage      State = "nextpage"
	StateClick         State = "click"
	StateAddToCart     State = "addtocart"
	StatePurchase      State = "purchase"
	StateExit          State = "exit" // terminal
	StateNone          State = ""     // ← 추가 (Back 처리용)
)

const (
	// Events (행위)
	EventSearchSubmitted EventType = "search_submitted"
	EventPageViewed      EventType = "page_viewed"
	EventPageClicked     EventType = "event_page_clicked"
	EventProductClicked  EventType = "product_clicked"
	EventCategoryClicked EventType = "category_clicked"
	EventAddToCart       EventType = "add_to_cart"
	EventPurchased       EventType = "purchased"
	EventBack            EventType = "back"
	EventExit            EventType = "exit"
)

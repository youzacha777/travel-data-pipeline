package event

type Event struct {
	EventID    string          `json:"event_id"`
	EventType  string          `json:"event_type"`
	EventTs    int64           `json:"event_ts"` // epoch millis
	UserID     string          `json:"user_id"`
	SessionID  string          `json:"session_id"`
	Attributes EventAttributes `json:"attributes"` // 유저 행동 구체 정보
}

type EventAttributes struct {
	State     string         `json:"state"` // FSM next state
	PrevState string         `json:"prev_state,omitempty"`
	Page      string         `json:"page,omitempty"` // browsing 페이지 정보
	Query     string         `json:"query,omitempty"`
	Product   *ProductInfo   `json:"product,omitempty"`
	Device    string         `json:"device,omitempty"`
	Referrer  string         `json:"referrer,omitempty"`
	Extra     map[string]any `json:"extra,omitempty"`
}

type ProductInfo struct {
	Country    string `json:"country,omitempty"`
	Category   string `json:"category,omitempty"`
	VendorType string `json:"vendor_type,omitempty"`
	ProductID  string `json:"product_id,omitempty"`
	Price      int    `json:"price,omitempty"`
}

package user

import (
	"event-generator/internal/fsm"
	"sync"
	"time"
)

// Session : 세션 기본 구조체
type Session struct {
	ID                      string
	UserID                  string
	State                   fsm.State
	PrevState               fsm.State
	PageType                string
	EventPage               string
	BrowsingCountryCategory string
	BrowsingProductCategory string
	SearchKeyword           string
	PageIndex               int
	LastEventTs             int64
	ExpiresAt               int64
	LastPicked              string
	LastProductID           string
	LastCategory            string
	LastCountry             string
	LastQuantity            int
}

// 전역 세션 저장소
var sessionStore = make(map[string]*Session)
var mu sync.Mutex // 동시성 문제 해결을 위한 뮤텍스

// NewSession creates a new Session for given userID with ttl duration.
func NewSession(sessionID, userID string, ttl time.Duration) *Session {
	now := time.Now().UnixMilli()
	return &Session{
		ID:          sessionID,
		UserID:      userID,
		State:       fsm.StateBrowsing, // 초기 상태
		LastEventTs: now,
		ExpiresAt:   time.Now().Add(ttl).UnixMilli(),
	}
}

// GetSession retrieves the session by sessionID, or creates a new one if not found
func GetSession(sessionID, userID string, ttl time.Duration) *Session {
	mu.Lock()
	defer mu.Unlock()

	// 세션이 이미 존재하면 반환
	if session, exists := sessionStore[sessionID]; exists {
		// 세션이 만료되었으면 새로 생성
		if time.Now().UnixMilli() > session.ExpiresAt {
			session = NewSession(sessionID, userID, ttl)
			sessionStore[sessionID] = session
		}
		return session
	}

	// 세션이 없으면 새로 생성
	session := NewSession(sessionID, userID, ttl)
	sessionStore[sessionID] = session
	return session
}

// 세션 종료 (삭제)
func DeleteSession(sessionID string) {
	mu.Lock()
	defer mu.Unlock()
	delete(sessionStore, sessionID)
}

// 세션 인터페이스 구현
// ===== identity =====
func (s *Session) GetID() string {
	return s.ID
}

func (s *Session) GetUserID() string {
	return s.UserID
}

// ===== state =====
func (s *Session) GetState() fsm.State {
	return s.State
}

func (s *Session) SetState(state fsm.State) {
	s.State = state
}

func (s *Session) GetPrevState() fsm.State {
	return s.PrevState
}

func (s *Session) SetPrevState(state fsm.State) {
	s.PrevState = state
}

// ===== time =====
func (s *Session) GetLastEventTs() int64 {
	return s.LastEventTs
}

func (s *Session) SetLastEventTs(ts int64) {
	s.LastEventTs = ts
}

// ===== browsing context =====
func (s *Session) SetPageType(pageType string) {
	s.PageType = pageType
}

func (s *Session) GetPageType() string {
	return s.PageType
}

func (s *Session) SetEventPage(eventPage string) {
	s.EventPage = eventPage
}

func (s *Session) GetEventPage() string {
	return s.EventPage
}

func (s *Session) SetBrowsingCountryCategory(country string) {
	s.BrowsingCountryCategory = country
}

func (s *Session) GetBrowsingCountryCategory() string {
	return s.BrowsingCountryCategory
}

func (s *Session) SetBrowsingProductCategory(category string) {
	s.BrowsingProductCategory = category
}

func (s *Session) GetBrowsingProductCategory() string {
	return s.BrowsingProductCategory
}

func (s *Session) ResetBrowsingContext() {
	s.BrowsingCountryCategory = ""
	s.BrowsingProductCategory = ""
}

// ===== search =====
func (s *Session) GetSearchKeyword() string {
	return s.SearchKeyword
}

func (s *Session) SetSearchKeyword(k string) {
	s.SearchKeyword = k
}

func (s *Session) SetPageIndex(i int) {
	s.PageIndex = i
}

func (s *Session) GetPageIndex() int {
	return s.PageIndex
}

func (s *Session) IncrementPageIndex() {
	s.PageIndex++
}

func (s *Session) SetExpiresAt(ts int64) {
	s.ExpiresAt = ts
}

func (s *Session) SetLastPicked(productID, category, country string) {
	s.LastProductID = productID
	s.LastCategory = category
	s.LastCountry = country
}

func (s *Session) GetLastPicked() (productID, category, country string) {
	return s.LastProductID, s.LastCategory, s.LastCountry
}

func (s *Session) GetLastCountry() string {
	return s.LastCountry
}

func (s *Session) GetLastProductID() string {
	return s.LastProductID
}

func (s *Session) GetLastCategory() string {
	return s.LastCategory
}

func (s *Session) SetLastQuantity(qty int) {
	s.LastQuantity = qty
}

func (s *Session) GetLastQuantity() int {
	return s.LastQuantity
}

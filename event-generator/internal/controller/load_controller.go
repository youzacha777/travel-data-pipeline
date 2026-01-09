package controller

import (
	"event-generator/internal/user"
	"fmt"
	"time"
)

type LoadController struct {
	TargetTPS int
	// PurchaseRatio float64

	UserPool       *user.UserPool
	SessionManager *user.SessionManager

	ticker       *time.Ticker
	quitChan     chan struct{}
	tickInterval time.Duration
}

// NewLoadController
func NewLoadController(
	tps int,
	// ratio float64,
	up *user.UserPool,
	sm *user.SessionManager,
) *LoadController {

	return &LoadController{
		TargetTPS: tps,
		// PurchaseRatio:  ratio,
		UserPool:       up,
		SessionManager: sm,
		quitChan:       make(chan struct{}),
		tickInterval:   10 * time.Millisecond, // ✅ 고정
	}
}

// Start begins load control loop
func (lc *LoadController) Start() {
	lc.ticker = time.NewTicker(lc.tickInterval)
	fmt.Printf(
		"[LoadController] started (TargetTPS=%d, tick=%s)\n",
		lc.TargetTPS,
		lc.tickInterval,
	)

	for {
		select {
		case <-lc.ticker.C:
			lc.tick()
		case <-lc.quitChan:
			lc.ticker.Stop()
			fmt.Println("[LoadController] stopped")
			return
		}
	}
}

// Stop stops the controller
func (lc *LoadController) Stop() {
	close(lc.quitChan)
}

// tick runs once per tickInterval
func (lc *LoadController) tick() {

	// ticks per second 계산
	ticksPerSecond := int(time.Second / lc.tickInterval)

	// tick 당 처리량
	batch := lc.TargetTPS / ticksPerSecond
	if batch <= 0 {
		batch = 1
	}

	// 유저 풀 확보
	lc.UserPool.EnsureUsers(lc.requiredUserCount())

	// FSM Step 실행 (이벤트 생성)
	for i := 0; i < batch; i++ {
		lc.SessionManager.Step()
	}
	// 고루틴을 사용하여 병렬로 이벤트 처리
	// for i := 0; i < batch; i++ {
	// 	go func() {
	// 		lc.SessionManager.Step() // 이벤트 생성
	// 	}()
	// }
}

// requiredUserCount calculates required users
func (lc *LoadController) requiredUserCount() int {
	// 경험적 기준 (조정 가능)
	return lc.TargetTPS * 2
}

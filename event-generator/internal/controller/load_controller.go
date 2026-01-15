package controller

import (
	"event-generator/internal/user"
	"fmt"
	"time"
)

type LoadController struct {
	TargetTPS int

	UserPool       *user.UserPool
	SessionManager *user.SessionManager

	ticker       *time.Ticker
	quitChan     chan struct{}
	tickInterval time.Duration

	workerCount int
}

func NewLoadController(
	tps int,
	up *user.UserPool,
	sm *user.SessionManager,
) *LoadController {
	// 2만 TPS 대응을 위해 워커 수를 CPU 코어 수의 2배 정도로 설정 권장
	// 예: 8코어 노트북이면 16개
	workerCount := 4

	return &LoadController{
		TargetTPS:      tps,
		UserPool:       up,
		SessionManager: sm,
		quitChan:       make(chan struct{}),
		tickInterval:   20 * time.Millisecond, // 10ms보다 20ms~50ms가 타이머 오차가 적고 안정적입니다.
		workerCount:    workerCount,
	}
}

func (lc *LoadController) Start() {
	// 1. 작업을 전달할 채널 (버퍼를 두어 송신자가 대기하지 않도록 함)
	taskCh := make(chan int, lc.workerCount*2)

	// 2. 워커 고루틴 풀 미리 생성 (딱 한 번만 실행됨)
	for w := 0; w < lc.workerCount; w++ {
		go func(id int) {
			for batchSize := range taskCh {
				for i := 0; i < batchSize; i++ {
					// 실제 이벤트 생성 로직 수행
					lc.SessionManager.Step()
				}
			}
		}(w)
	}

	lc.ticker = time.NewTicker(lc.tickInterval)
	defer lc.ticker.Stop()
	defer close(taskCh)

	fmt.Printf("[LoadController] started (TargetTPS=%d, tick=%s, workers=%d)\n",
		lc.TargetTPS, lc.tickInterval, lc.workerCount)

	ticksPerSecond := int(time.Second / lc.tickInterval)
	totalBatchPerTick := lc.TargetTPS / ticksPerSecond
	perWorkerBatch := totalBatchPerTick / lc.workerCount

	if perWorkerBatch <= 0 {
		perWorkerBatch = 1
	}

	for {
		select {
		case <-lc.ticker.C:
			// 유저 풀 확보
			lc.UserPool.EnsureUsers(lc.requiredUserCount())

			// 3. 고루틴 생성 없이 채널로 작업 지시만 내림 (매우 빠름)
			for w := 0; w < lc.workerCount; w++ {
				taskCh <- perWorkerBatch
			}

		case <-lc.quitChan:
			fmt.Println("[LoadController] stopping...")
			return
		}
	}
}

func (lc *LoadController) Stop() {
	close(lc.quitChan)
}

func (lc *LoadController) requiredUserCount() int {
	return lc.TargetTPS * 2
}

package main

import (
	"context"
	"event-generator/internal/controller"
	"event-generator/internal/event"
	"event-generator/internal/fsm"
	"event-generator/internal/generator"
	"event-generator/internal/metrics"
	"event-generator/internal/user"
	"event-generator/internal/worker"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	// 1. 모든 코어 활용 설정
	runtime.GOMAXPROCS(runtime.NumCPU())
	ctx, cancel := context.WithCancel(context.Background())

	// ======================
	// Metrics
	// ======================
	metricStore := metrics.NewInMemory()

	// ======================
	// Event Channel (버퍼 크기를 10만으로 유지)
	// ======================
	eventCh := make(chan *event.Event, 100000)

	// ======================
	// Core Components
	// ======================
	userPool := user.NewUserPool()
	userPool.EnsureUsers(100000)

	// [수정] 이제 main에서 전역 rand를 직접 시딩하거나 전달할 필요가 없습니다.
	// fsm과 generator 모두 내부적으로 math/rand/v2의 전역 소스를 사용합니다.
	fsmEngine := fsm.NewSimpleFSM()               // 인자 제거
	payloadGen := generator.NewPayloadGenerator() // 인자 없음 확인

	// ======================
	// Session Manager
	// ======================
	sm := user.NewSessionManager(
		userPool,
		fsmEngine,
		payloadGen,
		eventCh,
		metricStore,
		30*time.Minute,
	)

	// ======================
	// Load Controller
	// ======================
	targetTPS := 20000
	loadController := controller.NewLoadController(
		targetTPS,
		userPool,
		sm,
	)
	go loadController.Start()

	// ======================
	// Workers (Kafka Producer)
	// ======================
	// [성능 팁] TPS 2만 이상에서는 워커 수를 CPU 코어 수(runtime.NumCPU()) 정도로 늘리는 것이 유리합니다.
	workerCount := 12 // runtime.NumCPU()
	fmt.Printf("[MAIN] Using %d workers (CPU=%d)\n", workerCount, runtime.NumCPU())

	for i := 0; i < workerCount; i++ {
		w := worker.NewWorker(
			i,
			eventCh,
			metricStore,
			"localhost:9092",
			"user_events",
		)
		go w.Run(ctx)
	}

	// ======================
	// Metrics Snapshot & Channel Lag Monitor
	// ======================
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				snapshot := metricStore.Snapshot()
				fmt.Printf("[METRICS] %v | Lag: %d/%d\n",
					snapshot, len(eventCh), cap(eventCh))
			}
		}
	}()

	// ======================
	// Graceful Shutdown
	// ======================
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	fmt.Println("\n[MAIN] shutting down...")

	// 1. 이벤트 생성 중단
	loadController.Stop()

	// 2. 채널에 남은 이벤트 소비 대기
	fmt.Println("[MAIN] Draining event channel...")
	drainCtx, drainCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer drainCancel()

	for len(eventCh) > 0 {
		select {
		case <-drainCtx.Done():
			fmt.Println("[MAIN] Drain timeout - some events might be lost")
			goto ForceStop
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

ForceStop:
	// 3. 워커 종료 및 전송 플러시
	cancel()

	time.Sleep(2 * time.Second)
	fmt.Println("[MAIN] shutdown complete")
}

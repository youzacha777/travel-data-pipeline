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
	"math/rand"
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
	// Event Channel (버퍼 크기를 10만으로 늘림)
	// ======================
	// 초당 2만 개를 처리하므로 버퍼가 너무 작으면 생성부가 금방 막힙니다.
	eventCh := make(chan *event.Event, 100000)

	// ======================
	// Core Components
	// ======================
	userPool := user.NewUserPool()
	userPool.EnsureUsers(40000)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	fsmEngine := fsm.NewSimpleFSM(rnd)
	payloadGen := generator.NewPayloadGenerator()

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
	// 우선 4개로 고정하셨으니 4개로 둡니다.
	workerCount := 4
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
		// 실시간 확인을 위해 1초마다 출력합니다.
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				snapshot := metricStore.Snapshot()
				// len(eventCh)가 현재 채널에 쌓여있는 대기열(Lag)입니다.
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

	// 2. 채널에 남은 이벤트가 소비될 때까지 잠시 대기
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

	// Kafka 비동기 전송 처리를 위해 약간 더 대기
	time.Sleep(2 * time.Second)
	fmt.Println("[MAIN] shutdown complete")
}

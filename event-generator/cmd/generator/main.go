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
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ======================
	// Metrics
	// ======================
	metricStore := metrics.NewInMemory()

	// ======================
	// Event Channel
	// ======================
	eventCh := make(chan *event.Event, 100)

	// ======================
	// Core Components
	// ======================
	userPool := user.NewUserPool()
	userPool.EnsureUsers(20000)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	fsmEngine := fsm.NewSimpleFSM(rnd)

	payloadGen := generator.NewPayloadGenerator()

	// ======================
	// Session Manager
	// ======================
	sm := user.NewSessionManager(
		userPool,
		fsmEngine,
		payloadGen, // 인터페이스 타입으로 들어감
		eventCh,
		metricStore,
		30*time.Minute,
	)

	// ======================
	// Load Controller (TPS 제어의 유일한 진입점)
	// ======================
	loadController := controller.NewLoadController(
		10000, // 🔥 여기서 TPS 조절 (10k / 50k)
		// 0.2,
		userPool,
		sm,
	)

	go loadController.Start()

	// ======================
	// Workers (Event Consumer)
	// ======================
	workerCount := 4
	for i := 0; i < workerCount; i++ {
		w := worker.NewWorker(
			i,
			eventCh,
			metricStore,
			"localhost:9092", // Kafka 브로커 주소
			"user_events",    // Kafka 토픽
		)
		go w.Run(ctx)
	}

	// ======================
	// Metrics Snapshot
	// ======================
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fmt.Printf("[METRICS] %+v\n", metricStore.Snapshot())
			}
		}
	}()

	// ======================
	// Shutdown
	// ======================
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	fmt.Println("shutting down...")
	cancel()
	loadController.Stop()
	time.Sleep(time.Second)
}

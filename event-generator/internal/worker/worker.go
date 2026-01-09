package worker

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"event-generator/internal/event"
	"event-generator/internal/metrics"

	"github.com/segmentio/kafka-go"
)

type Worker struct {
	id        int
	eventCh   <-chan *event.Event
	metrics   metrics.Metrics
	kafkaAddr string
	topic     string
}

func NewWorker(
	id int,
	eventCh <-chan *event.Event,
	m metrics.Metrics,
	kafkaAddr, topic string,
) *Worker {
	return &Worker{
		id:        id,
		eventCh:   eventCh,
		metrics:   m,
		kafkaAddr: kafkaAddr,
		topic:     topic,
	}
}

func (w *Worker) Run(ctx context.Context) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{w.kafkaAddr},
		Topic:    w.topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	// 메시지를 한번에 처리할 수 있도록 버퍼 사용
	var messageBatch []kafka.Message
	const batchSize = 100 // 배치 크기 (100개 메시지마다 전송)

	var mu sync.Mutex // 동기화를 위한 뮤텍스

	for {
		select {
		case <-ctx.Done():
			return

		case ev := <-w.eventCh:
			// 이벤트를 Kafka 메시지로 변환
			msgBytes, err := json.Marshal(ev)
			if err != nil {
				log.Println("Event marshal error:", err)
				continue
			}

			msg := kafka.Message{
				Key:   []byte(ev.UserID),
				Value: msgBytes,
			}

			// 배치 크기마다 메시지 전송
			mu.Lock()
			messageBatch = append(messageBatch, msg)
			if len(messageBatch) >= batchSize {
				err = writer.WriteMessages(ctx, messageBatch...)
				if err != nil {
					log.Println("Kafka batch write error:", err)
				}
				messageBatch = messageBatch[:0] // 배치 초기화
			}
			mu.Unlock()

			// 메트릭 증가
			w.metrics.IncEvent(ev.EventType)

		}
	}
}

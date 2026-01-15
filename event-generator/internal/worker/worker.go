package worker

import (
	"context"
	"encoding/json"
	"time"

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
	writer := &kafka.Writer{
		Addr:     kafka.TCP(w.kafkaAddr),
		Topic:    w.topic,
		Balancer: &kafka.Hash{},
		// 수동 배칭 대신 라이브러리 설정을 활용
		BatchSize:    1000,
		BatchTimeout: 50 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		Async:        true, // 비동기 모드이므로 WriteMessages는 논블로킹
		Compression:  kafka.Snappy,
	}
	defer writer.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-w.eventCh:
			if !ok {
				return
			}

			// 1. JSON 직렬화
			msgBytes, err := json.Marshal(ev)
			if err != nil {
				w.metrics.IncError("marshal_event")
				continue
			}

			// 2. 즉시 쓰기 (Async 모드라 내부 버퍼로 바로 들어감)
			err = writer.WriteMessages(ctx, kafka.Message{
				Key:   []byte(ev.UserID),
				Value: msgBytes,
			})

			if err != nil {
				w.metrics.IncError("write_messages")
			} else {
				// 3. 성공 시 메트릭 업데이트
				w.metrics.IncEvent(ev.EventType)
			}

			// TIP: 2만 TPS 이상을 원하신다면 여기서 한 번에 여러 개를
			// 꺼내는 for 루프를 추가할 수 있지만, Async 모드에선 이정도로도 충분합니다.
		}
	}
}

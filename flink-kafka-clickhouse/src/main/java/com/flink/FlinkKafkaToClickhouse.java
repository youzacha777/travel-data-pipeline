package com.flink;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.api.common.serialization.SimpleStringSchema;
import org.apache.flink.connector.jdbc.JdbcConnectionOptions;
import org.apache.flink.connector.jdbc.JdbcExecutionOptions;
import org.apache.flink.connector.jdbc.JdbcSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.api.common.ExecutionConfig;
import org.apache.flink.api.common.restartstrategy.RestartStrategies;
import org.apache.flink.runtime.state.filesystem.FsStateBackend;
import org.apache.flink.streaming.api.CheckpointingMode;

public class FlinkKafkaToClickhouse {

    private static final ObjectMapper MAPPER = new ObjectMapper();

    public static void main(String[] args) throws Exception {

        // 1. Flink 환경 설정
        final StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        env.setParallelism(4);  // 병렬 처리 수 설정

        // 1.1 ExecutionConfig 설정
        ExecutionConfig config = env.getConfig();
        config.setAutoWatermarkInterval(2000); // 수분 마커 자동 생성 주기 2초
        config.setTaskCancellationInterval(10000); // 작업 취소 간격 10초

        // 1.2 체크포인트 설정
        env.enableCheckpointing(30000, CheckpointingMode.EXACTLY_ONCE); // 30초마다 체크포인트
        env.getCheckpointConfig().setMinPauseBetweenCheckpoints(5000); // 최소 대기 시간
        env.setStateBackend(new FsStateBackend("file:///tmp/flink-checkpoints")); // 체크포인트 상태 저장 위치

        // 1.3 재시도 전략 설정
        env.setRestartStrategy(RestartStrategies.fixedDelayRestart(3, 10000)); // 최대 3회 재시도, 10초 대기

        // 2. Kafka Source 설정 (OffsetsInitializer 추가 권장)
        KafkaSource<String> kafkaSource = KafkaSource.<String>builder()
                .setBootstrapServers("kafka:29092")
                .setTopics("user_events")
                .setGroupId("flink-user-events")
                .setStartingOffsets(OffsetsInitializer.earliest()) // 시작점 명시
                .setValueOnlyDeserializer(new SimpleStringSchema())
                .build();

        // Kafka에서 데이터를 읽어오는 스트림
        DataStream<String> stream = env.fromSource(kafkaSource, WatermarkStrategy.noWatermarks(), "Kafka Source");

        // 3. ClickHouse JDBC 설정
        JdbcConnectionOptions jdbcOptions = new JdbcConnectionOptions.JdbcConnectionOptionsBuilder()
                .withUrl("jdbc:clickhouse://clickhouse:8123/user_events") // 공식 드라이버에 맞는 URL
                .withDriverName("com.clickhouse.jdbc.ClickHouseDriver") // 최신 ClickHouse 드라이버 사용
                .withUsername("clickhouse")
                .withPassword("pass")
                .build();

        JdbcExecutionOptions executionOptions = JdbcExecutionOptions.builder()
                .withBatchSize(500)         // 배치 크기
                .withBatchIntervalMs(2000)  // 2초마다 flush
                .withMaxRetries(3)          // 실패 시 3회 재시도
                .build();

        // 4. Sink 설정 (ClickHouse에 데이터 삽입)
        stream.addSink(
                JdbcSink.sink(
                        "INSERT INTO user_events_raw (event_id, user_id, session_id, event_type, event_ts, state, payload) VALUES (?, ?, ?, ?, ?, ?, ?)",
                        (ps, value) -> {
                            try {
                                JsonNode json = MAPPER.readTree(value);
                                ps.setString(1, json.path("event_id").asText());
                                ps.setString(2, json.path("user_id").asText());
                                ps.setString(3, json.path("session_id").asText());
                                ps.setString(4, json.path("event_type").asText());
                                ps.setLong(5, json.path("event_ts").asLong());
                                ps.setString(6, json.path("attributes").path("state").asText());
                                ps.setString(7, value); // payload 전체 JSON 저장
                            } catch (Exception e) {
                                // 에러 로깅 시 로깅 프레임워크 사용 권장
                                System.err.println("JSON Parsing Error: " + value);
                            }
                        },
                        executionOptions,
                        jdbcOptions
                )
        ).name("ClickHouse Sink");

        // 5. 플링크 작업 실행
        env.execute("Flink Kafka → ClickHouse Pipeline");
    }
}

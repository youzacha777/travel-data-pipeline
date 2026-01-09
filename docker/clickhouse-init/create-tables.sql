-- DB 생성
CREATE DATABASE IF NOT EXISTS user_events;

-- 테이블 생성
CREATE TABLE IF NOT EXISTS user_events.user_events_raw (
    event_id String,
    user_id String,
    session_id String,
    event_type String,
    event_ts UInt64,
    state String,
    payload String
) ENGINE = MergeTree()
ORDER BY (session_id, event_ts);


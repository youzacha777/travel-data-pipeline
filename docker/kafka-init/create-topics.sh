#!/bin/bash
set -e

TOPIC_NAME="user_events"
BOOTSTRAP_SERVER="kafka:29092" # kafka-init 컨테이너에서도 kafka 호스트명으로 접속 가능
PARTITIONS=16      # 초당 20k 이벤트 목표
REPLICATION=1     # 단일 브로커 환경

echo "Waiting for Kafka broker to be ready..."

# Kafka 준비될 때까지 대기
until kafka-topics --bootstrap-server $BOOTSTRAP_SERVER --list > /dev/null 2>&1; do
  echo "Kafka not ready yet, sleeping 3s..."
  sleep 3
done

echo "Kafka is ready! Checking if topic '$TOPIC_NAME' exists..."

# 토픽 존재 여부 확인 후 생성
if kafka-topics --bootstrap-server $BOOTSTRAP_SERVER --list | grep -q "^${TOPIC_NAME}$"; then
    echo "Topic '$TOPIC_NAME' already exists, skipping creation."
else
    echo "Creating topic '$TOPIC_NAME'..."
    kafka-topics --create \
        --bootstrap-server $BOOTSTRAP_SERVER \
        --replication-factor $REPLICATION \
        --partitions $PARTITIONS \
        --topic $TOPIC_NAME
    echo "Topic '$TOPIC_NAME' created successfully!"
fi

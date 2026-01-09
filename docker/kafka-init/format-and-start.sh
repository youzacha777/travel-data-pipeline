#!/bin/bash
set -e

# Kafka 데이터 디렉토리
KAFKA_DATA_DIR="/var/lib/kafka/data"

# CLUSTER_ID 설정 (원하면 환경변수로도 받을 수 있음)
CLUSTER_ID="${KAFKA_CLUSTER_ID:-my-cluster-1}"

# Kafka 포맷 (이미 포맷되어 있으면 --ignore-formatted 사용)
if [ ! -d "$KAFKA_DATA_DIR/meta.properties" ]; then
    echo "Formatting Kafka storage..."
    kafka-storage format \
        --ignore-formatted \
        --cluster-id $CLUSTER_ID \
        --config /etc/kafka/kraft/server.properties
else
    echo "Kafka storage already formatted, skipping..."
fi

# Kafka 실행 (Confluent 기본 entrypoint 실행)
echo "Starting Kafka broker..."
/etc/confluent/docker/run

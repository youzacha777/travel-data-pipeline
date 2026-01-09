#!/bin/sh
set -e

# -----------------------------
# 환경 변수
# -----------------------------
FLINK_REST="http://flink-jobmanager:8081"
JOB_JAR="/opt/flink/usrlib/flink-kafka-clickhouse-1.0-SNAPSHOT.jar"
LOG_FILE="/opt/flink/usrlib/flink-job-init.log"

echo "[INIT] Starting Flink job submission..." | tee -a ${LOG_FILE}

# -----------------------------
# 1. Flink REST API 준비 대기
# -----------------------------
echo "[INIT] Waiting for Flink REST API..." | tee -a ${LOG_FILE}
until curl -sf ${FLINK_REST}/jobs/overview > /dev/null; do
  sleep 3
done
echo "[INIT] Flink REST API is ready." | tee -a ${LOG_FILE}

# -----------------------------
# 2. 이미 실행 중인 Job 확인
# -----------------------------
RUNNING_JOBS=$(curl -s ${FLINK_REST}/jobs/overview | grep RUNNING || true)
if [ -n "$RUNNING_JOBS" ]; then
  echo "[INIT] Flink job already running. Skip submission." | tee -a ${LOG_FILE}
  exit 0
fi

# -----------------------------
# 3. Flink Job 제출
# -----------------------------
echo "[INIT] Submitting Flink job..." | tee -a ${LOG_FILE}
flink run -d \
  -m flink-jobmanager:8081 \
  -c com.flink.FlinkKafkaToClickhouse \
  ${JOB_JAR} 2>&1 | tee -a ${LOG_FILE}

# -----------------------------
# 4. Job 제출 확인
# -----------------------------
JOB_ID=$(flink list -r | grep "FlinkKafkaToClickhouse" | awk '{print $1}' || true)
if [ -z "$JOB_ID" ]; then
  echo "[ERROR] Job submission failed!" | tee -a ${LOG_FILE}
  exit 1
else
  echo "[INIT] Job submitted successfully. Job ID: $JOB_ID" | tee -a ${LOG_FILE}
fi

#!/usr/bin/env bash

export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-test}"
export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-test}"
export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-${AWS_REGION:-us-east-1}}"

AWS_ENDPOINT_URL="${AWS_ENDPOINT_URL:-${ENDPOINT_URL:-http://localhost:4566}}"
RAW_QUEUE_URL="${RAW_QUEUE_URL:-http://localhost:4566/000000000000/raw-events}"
PROCESSED_QUEUE_URL="${PROCESSED_QUEUE_URL:-http://localhost:4566/000000000000/processed-events}"
RAW_DLQ_URL="${RAW_DLQ_URL:-http://localhost:4566/000000000000/raw-events-dlq}"
API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
SMOKE_TIMEOUT_SECONDS="${SMOKE_TIMEOUT_SECONDS:-45}"

require_command() {
  local command_name="$1"
  if ! command -v "$command_name" >/dev/null 2>&1; then
    echo "Missing required command: $command_name" >&2
    exit 1
  fi
}

require_smoke_commands() {
  require_command aws
  require_command curl
  require_command python3
}

uuid_v4() {
  python3 -c 'import uuid; print(uuid.uuid4())'
}

send_raw_message() {
  local body="$1"
  aws --endpoint-url="$AWS_ENDPOINT_URL" sqs send-message \
    --queue-url "$RAW_QUEUE_URL" \
    --region "$AWS_DEFAULT_REGION" \
    --message-body "$body" >/dev/null
}

receive_raw_dlq_messages() {
  aws --endpoint-url="$AWS_ENDPOINT_URL" sqs receive-message \
    --queue-url "$RAW_DLQ_URL" \
    --region "$AWS_DEFAULT_REGION" \
    --max-number-of-messages 10 \
    --wait-time-seconds 1 \
    --attribute-names ApproximateReceiveCount \
    --message-attribute-names All
}

purge_queue() {
  local queue_url="$1"
  aws --endpoint-url="$AWS_ENDPOINT_URL" sqs purge-queue \
    --queue-url "$queue_url" \
    --region "$AWS_DEFAULT_REGION" >/dev/null 2>&1 || true
}

wait_for_health() {
  local deadline=$((SECONDS + SMOKE_TIMEOUT_SECONDS))
  until curl -fsS "$API_BASE_URL/health" >/dev/null 2>&1; do
    if (( SECONDS >= deadline )); then
      echo "API health did not become ready within ${SMOKE_TIMEOUT_SECONDS}s" >&2
      return 1
    fi
    sleep 1
  done
}

get_summary() {
  local developer_id="$1"
  curl -fsS "$API_BASE_URL/metrics/${developer_id}/summary"
}

get_metrics() {
  local developer_id="$1"
  curl -fsS "$API_BASE_URL/metrics/${developer_id}"
}

json_field() {
  local json="$1"
  local field="$2"
  python3 -c 'import json, sys
data = json.loads(sys.argv[1])
value = data
for part in sys.argv[2].split("."):
    value = value[part]
print(value)' "$json" "$field"
}

assert_json_field_equals() {
  local json="$1"
  local field="$2"
  local expected="$3"
  local actual
  actual="$(json_field "$json" "$field")"
  if [[ "$actual" != "$expected" ]]; then
    echo "Assertion failed for ${field}: expected ${expected}, got ${actual}" >&2
    echo "$json" >&2
    exit 1
  fi
}

wait_for_summary_field() {
  local developer_id="$1"
  local field="$2"
  local expected="$3"
  local deadline=$((SECONDS + SMOKE_TIMEOUT_SECONDS))
  local summary=""
  local actual=""

  while (( SECONDS < deadline )); do
    if summary="$(get_summary "$developer_id" 2>/dev/null)"; then
      actual="$(json_field "$summary" "$field" 2>/dev/null || true)"
      if [[ "$actual" == "$expected" ]]; then
        echo "$summary"
        return 0
      fi
    fi
    sleep 1
  done

  echo "Timed out waiting for ${developer_id}.${field}=${expected}; last value=${actual}" >&2
  [[ -n "$summary" ]] && echo "$summary" >&2
  return 1
}

assert_metrics_count() {
  local developer_id="$1"
  local expected="$2"
  local metrics
  local actual
  metrics="$(get_metrics "$developer_id")"
  actual="$(python3 -c 'import json, sys; print(len(json.loads(sys.argv[1])))' "$metrics")"
  if [[ "$actual" != "$expected" ]]; then
    echo "Assertion failed for metrics count: expected ${expected}, got ${actual}" >&2
    echo "$metrics" >&2
    exit 1
  fi
}

wait_for_raw_dlq_event() {
  local event_id="$1"
  local deadline=$((SECONDS + SMOKE_TIMEOUT_SECONDS))
  local messages=""

  while (( SECONDS < deadline )); do
    messages="$(receive_raw_dlq_messages || true)"
    if python3 -c 'import json, sys
event_id = sys.argv[1]
raw = sys.argv[2].strip()
if not raw:
    sys.exit(1)
data = json.loads(raw)
for message in data.get("Messages", []):
    body = json.loads(message.get("Body", "{}"))
    if body.get("event_id") == event_id:
        sys.exit(0)
sys.exit(1)' "$event_id" "$messages"; then
      echo "$messages"
      return 0
    fi
    sleep 2
  done

  echo "Timed out waiting for invalid event ${event_id} in raw-events-dlq" >&2
  return 1
}

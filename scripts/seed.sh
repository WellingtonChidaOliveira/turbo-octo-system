#!/usr/bin/env bash
set -euo pipefail

export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-test}"
export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-test}"
export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-us-east-1}"

ENDPOINT_URL="${ENDPOINT_URL:-http://localhost:4566}"
QUEUE_URL="${RAW_QUEUE_URL:-http://localhost:4566/000000000000/raw-events}"

send_message() {
  local body="$1"
  aws --endpoint-url="$ENDPOINT_URL" sqs send-message \
    --queue-url "$QUEUE_URL" \
    --region us-east-1 \
    --message-body "$body" >/dev/null
}

echo "Seeding raw-events queue at $QUEUE_URL"

send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440000","developer_id":"dev-001","metric_type":"commits","value":12,"repository":"org/api","timestamp":"2026-04-15T10:30:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440001","developer_id":"dev-001","metric_type":"pull_requests","value":3,"repository":"org/api","timestamp":"2026-04-15T10:31:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440002","developer_id":"dev-001","metric_type":"review_time_minutes","value":45,"repository":"org/api","timestamp":"2026-04-15T10:32:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440003","developer_id":"dev-002","metric_type":"commits","value":7,"repository":"org/web","timestamp":"2026-04-15T10:33:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440004","developer_id":"dev-002","metric_type":"pull_requests","value":2,"repository":"org/web","timestamp":"2026-04-15T10:34:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440005","developer_id":"dev-002","metric_type":"review_time_minutes","value":60,"repository":"org/web","timestamp":"2026-04-15T10:35:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440006","developer_id":"dev-003","metric_type":"commits","value":18,"repository":"org/mobile","timestamp":"2026-04-15T10:36:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440007","developer_id":"dev-003","metric_type":"pull_requests","value":5,"repository":"org/mobile","timestamp":"2026-04-15T10:37:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440008","developer_id":"dev-003","metric_type":"review_time_minutes","value":30,"repository":"org/mobile","timestamp":"2026-04-15T10:38:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440009","developer_id":"dev-004","metric_type":"commits","value":4,"repository":"org/worker","timestamp":"2026-04-15T10:39:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440010","developer_id":"dev-004","metric_type":"pull_requests","value":1,"repository":"org/worker","timestamp":"2026-04-15T10:40:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440011","developer_id":"dev-004","metric_type":"review_time_minutes","value":90,"repository":"org/worker","timestamp":"2026-04-15T10:41:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440012","developer_id":"dev-005","metric_type":"commits","value":22,"repository":"org/platform","timestamp":"2026-04-15T10:42:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440013","developer_id":"dev-005","metric_type":"pull_requests","value":6,"repository":"org/platform","timestamp":"2026-04-15T10:43:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440014","developer_id":"dev-005","metric_type":"review_time_minutes","value":120,"repository":"org/platform","timestamp":"2026-04-15T10:44:00Z"}'

# Invalid messages for validation/DLQ demonstration.
send_message '{"event_id":"invalid-uuid","developer_id":"dev-006","metric_type":"commits","value":1,"repository":"org/api","timestamp":"2026-04-15T10:45:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440016","developer_id":"","metric_type":"commits","value":1,"repository":"org/api","timestamp":"2026-04-15T10:46:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440017","developer_id":"dev-006","metric_type":"review_time_minutes","value":1500,"repository":"org/api","timestamp":"2026-04-15T10:47:00Z"}'

# Duplicates for the Aggregator idempotency stage.
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440018","developer_id":"dev-007","metric_type":"commits","value":9,"repository":"org/api","timestamp":"2026-04-15T10:48:00Z"}'
send_message '{"event_id":"550e8400-e29b-41d4-a716-446655440018","developer_id":"dev-007","metric_type":"commits","value":9,"repository":"org/api","timestamp":"2026-04-15T10:48:00Z"}'

echo "Seed complete"

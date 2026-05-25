#!/usr/bin/env bash
set -euo pipefail

export SMOKE_TIMEOUT_SECONDS="${SMOKE_TIMEOUT_SECONDS:-140}"

source "$(dirname "$0")/lib-smoke.sh"

require_smoke_commands
wait_for_health
purge_queue "$RAW_DLQ_URL"

developer_id="smoke-invalid-$(date +%s)"
event_id="$(uuid_v4)"

echo "Sending invalid raw event for ${developer_id}"
send_raw_message "{\"event_id\":\"${event_id}\",\"developer_id\":\"${developer_id}\",\"metric_type\":\"review_time_minutes\",\"value\":1441,\"repository\":\"org/smoke\",\"timestamp\":\"2026-04-15T11:20:00Z\"}"

wait_for_raw_dlq_event "$event_id" >/dev/null

if get_summary "$developer_id" >/dev/null 2>&1; then
  echo "Invalid event unexpectedly produced a summary for ${developer_id}" >&2
  exit 1
fi

echo "Smoke invalid passed"

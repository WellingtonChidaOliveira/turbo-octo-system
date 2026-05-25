#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/lib-smoke.sh"

require_smoke_commands
wait_for_health

developer_id="smoke-receive-$(date +%s)"
event_id="$(uuid_v4)"

echo "Sending one valid raw event for ${developer_id}"
send_raw_message "{\"event_id\":\"${event_id}\",\"developer_id\":\"${developer_id}\",\"metric_type\":\"commits\",\"value\":5,\"repository\":\"org/smoke\",\"timestamp\":\"2026-04-15T11:00:00Z\"}"

summary="$(wait_for_summary_field "$developer_id" "events_processed" "1")"
assert_json_field_equals "$summary" "total_commits" "5"
assert_json_field_equals "$summary" "last_activity" "2026-04-15T11:00:00Z"
assert_metrics_count "$developer_id" "1"

echo "Smoke receive passed"

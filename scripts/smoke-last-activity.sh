#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/lib-smoke.sh"

require_smoke_commands
wait_for_health

developer_id="smoke-last-activity-$(date +%s)"
newer_event_id="$(uuid_v4)"
older_event_id="$(uuid_v4)"
newer_timestamp="2026-04-15T12:10:00Z"
older_timestamp="2026-04-15T12:00:00Z"

echo "Sending newer event first for ${developer_id}"
send_raw_message "{\"event_id\":\"${newer_event_id}\",\"developer_id\":\"${developer_id}\",\"metric_type\":\"commits\",\"value\":3,\"repository\":\"org/smoke\",\"timestamp\":\"${newer_timestamp}\"}"

summary="$(wait_for_summary_field "$developer_id" "last_activity" "$newer_timestamp")"
assert_json_field_equals "$summary" "events_processed" "1"

echo "Sending older event second for ${developer_id}"
send_raw_message "{\"event_id\":\"${older_event_id}\",\"developer_id\":\"${developer_id}\",\"metric_type\":\"pull_requests\",\"value\":2,\"repository\":\"org/smoke\",\"timestamp\":\"${older_timestamp}\"}"

summary="$(wait_for_summary_field "$developer_id" "events_processed" "2")"
assert_json_field_equals "$summary" "last_activity" "$newer_timestamp"
assert_json_field_equals "$summary" "total_commits" "3"
assert_json_field_equals "$summary" "total_pull_requests" "2"
assert_metrics_count "$developer_id" "2"

echo "Smoke last-activity passed"

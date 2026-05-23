#!/usr/bin/env bash
set -euo pipefail

export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-test}"
export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-test}"
export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-us-east-1}"

ENDPOINT_URL="${ENDPOINT_URL:-http://localhost:4566}"
QUEUE_URL="${RAW_QUEUE_URL:-http://localhost:4566/000000000000/processed-events}"


aws sqs purge-queue --queue-url "$QUEUE_URL" --endpoint-url="$ENDPOINT_URL" --region us-east-1

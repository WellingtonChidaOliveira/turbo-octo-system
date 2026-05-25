#!/usr/bin/env bash
set -euo pipefail

export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-test}"
export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-test}"
export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-${AWS_REGION:-us-east-1}}"

AWS_ENDPOINT_URL="${AWS_ENDPOINT_URL:-${ENDPOINT_URL:-http://localhost:4566}}"
PROCESSED_QUEUE_URL="${PROCESSED_QUEUE_URL:-http://localhost:4566/000000000000/processed-events}"


aws sqs purge-queue --queue-url "$PROCESSED_QUEUE_URL" --endpoint-url="$AWS_ENDPOINT_URL" --region "$AWS_DEFAULT_REGION"

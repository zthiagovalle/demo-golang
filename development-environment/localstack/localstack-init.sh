#!/usr/bin/env bash
set -euo pipefail

REGION=us-east-1
TOPIC_NAME=product-events
QUEUE_NAME=product-created-queue

echo "[localstack-init] creating SNS topic: $TOPIC_NAME"
TOPIC_ARN=$(awslocal sns create-topic --name "$TOPIC_NAME" --region "$REGION" --query TopicArn --output text)

echo "[localstack-init] creating SQS queue: $QUEUE_NAME"
QUEUE_URL=$(awslocal sqs create-queue --queue-name "$QUEUE_NAME" --region "$REGION" --query QueueUrl --output text)
QUEUE_ARN=$(awslocal sqs get-queue-attributes --queue-url "$QUEUE_URL" --attribute-names QueueArn --region "$REGION" --query 'Attributes.QueueArn' --output text)

echo "[localstack-init] subscribing $QUEUE_ARN to $TOPIC_ARN (raw delivery)"
awslocal sns subscribe \
  --topic-arn "$TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$QUEUE_ARN" \
  --attributes '{"RawMessageDelivery":"true"}' \
  --region "$REGION"

echo "[localstack-init] done"

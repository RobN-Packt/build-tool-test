#!/usr/bin/env bash
set -euo pipefail

AWS_ACCOUNT_ID=${AWS_ACCOUNT_ID:-123456789012}
AWS_REGION=${AWS_REGION:-eu-west-1}

docker build -t api:latest apps/api
docker tag api:latest "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/poc/api:latest"
docker push "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/poc/api:latest"

docker build -t web:latest apps/web
docker tag web:latest "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/poc/web:latest"
docker push "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/poc/web:latest"

docker build -t worker:latest apps/worker
docker tag worker:latest "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/poc/worker:latest"
docker push "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/poc/worker:latest"

aws ecs update-service --cluster poc-cluster --service poc-api --force-new-deployment
aws lambda update-function-code --function-name poc-purchase-worker --image-uri "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/poc/worker:latest"
aws s3 sync apps/web/.next/static s3://book-poc-web-static/_next/static --delete
if [ -n "${CLOUDFRONT_DISTRIBUTION_ID:-}" ]; then
  aws cloudfront create-invalidation --distribution-id "$CLOUDFRONT_DISTRIBUTION_ID" --paths "/*"
fi

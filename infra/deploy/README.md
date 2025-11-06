# Deployment Guide

This directory captures high-level workflows for deploying the PoC via the supported CI systems.

## Terraform
1. Authenticate with AWS (OIDC in CI, `aws sso login` locally)
2. `cd infra/terraform`
3. `terraform init`
4. `terraform plan -out tfplan`
5. `terraform apply tfplan`

The stack provisions ECR repositories, an ECS Fargate service for the API, a Lambda function for purchase processing, an SQS queue, and CloudFront + S3 for the web frontend.

## Docker Images
- API: `docker build -t $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/poc/api:TAG apps/api`
- Web: `docker build -t $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/poc/web:TAG apps/web`
- Worker: `docker build -t $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/poc/worker:TAG apps/worker`

Push after authenticating with `aws ecr get-login-password`.

## ECS Deploy
The GitHub Actions `deploy.yml` applies Terraform, pushes images, and triggers a new task set via `aws ecs update-service --force-new-deployment`.

## Lambda Deploy
The worker Docker image is published to ECR and referenced by the Lambda function. Updating the image tag and running `terraform apply` (or the CI workflow) publishes the new version.

## Frontend Deploy
Static assets are uploaded to the provisioned S3 bucket via the CI workflow, followed by a CloudFront invalidation to refresh caches.

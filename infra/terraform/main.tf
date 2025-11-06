terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

locals {
  project    = "book-poc"
  ecr_prefix = "poc"
}

resource "aws_ecr_repository" "api" {
  name = "${local.ecr_prefix}/api"
  image_scanning_configuration {
    scan_on_push = true
  }
  tags = var.tags
}

resource "aws_ecr_repository" "web" {
  name = "${local.ecr_prefix}/web"
  image_scanning_configuration {
    scan_on_push = true
  }
  tags = var.tags
}

resource "aws_ecr_repository" "worker" {
  name = "${local.ecr_prefix}/worker"
  image_scanning_configuration {
    scan_on_push = true
  }
  tags = var.tags
}

resource "aws_sqs_queue" "purchases" {
  name                       = "${local.project}-purchases"
  visibility_timeout_seconds = 60
  message_retention_seconds  = 86400
  tags                       = var.tags
}

module "ecs" {
  source         = "./ecs"
  aws_region     = var.aws_region
  cluster_name   = "poc-cluster"
  service_name   = "poc-api"
  container_image = "${var.aws_account_id}.dkr.ecr.${var.aws_region}.amazonaws.com/${aws_ecr_repository.api.name}:latest"
  queue_arn      = aws_sqs_queue.purchases.arn
  queue_url      = aws_sqs_queue.purchases.id
  tags           = var.tags
}

module "lambda" {
  source             = "./lambda"
  function_name      = "poc-purchase-worker"
  image_uri          = "${var.aws_account_id}.dkr.ecr.${var.aws_region}.amazonaws.com/${aws_ecr_repository.worker.name}:latest"
  queue_arn          = aws_sqs_queue.purchases.arn
  queue_url          = aws_sqs_queue.purchases.id
  environment        = {
    AWS_REGION = var.aws_region
  }
  tags               = var.tags
}

resource "aws_s3_bucket" "web" {
  bucket = "${local.project}-web-static"
  tags   = var.tags
}

resource "aws_s3_bucket_public_access_block" "web" {
  bucket                  = aws_s3_bucket.web.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_cloudfront_origin_access_control" "web" {
  name                              = "${local.project}-oac"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "web" {
  enabled             = true
  default_root_object = "index.html"

  origin {
    domain_name              = aws_s3_bucket.web.bucket_regional_domain_name
    origin_id                = "web-origin"
    origin_access_control_id = aws_cloudfront_origin_access_control.web.id
  }

  default_cache_behavior {
    target_origin_id       = "web-origin"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = var.tags
}

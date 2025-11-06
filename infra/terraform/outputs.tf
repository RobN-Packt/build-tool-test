output "ecr_api_repository" {
  value = aws_ecr_repository.api.repository_url
}

output "ecr_web_repository" {
  value = aws_ecr_repository.web.repository_url
}

output "ecr_worker_repository" {
  value = aws_ecr_repository.worker.repository_url
}

output "ecs_cluster" {
  value = module.ecs.cluster_id
}

output "ecs_service" {
  value = module.ecs.service_name
}

output "lambda_function" {
  value = module.lambda.function_name
}

output "cloudfront_domain" {
  value = aws_cloudfront_distribution.web.domain_name
}

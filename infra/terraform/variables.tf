variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "eu-west-1"
}

variable "aws_account_id" {
  description = "AWS account ID"
  type        = string
  default     = "123456789012"
}

variable "tags" {
  description = "Common tags for all resources"
  type        = map(string)
  default     = {
    Project = "book-poc"
  }
}

variable "web_domain" {
  description = "Optional custom domain for the web frontend"
  type        = string
  default     = ""
}

variable "aws_region" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "service_name" {
  type = string
}

variable "container_image" {
  type = string
}

variable "queue_arn" {
  type = string
}

variable "queue_url" {
  type = string
}

variable "tags" {
  type = map(string)
  default = {}
}

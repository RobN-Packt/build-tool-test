variable "function_name" {
  type = string
}

variable "image_uri" {
  type = string
}

variable "queue_arn" {
  type = string
}

variable "queue_url" {
  type = string
}

variable "environment" {
  type    = map(string)
  default = {}
}

variable "tags" {
  type    = map(string)
  default = {}
}

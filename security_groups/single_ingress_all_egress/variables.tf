variable "name" {
  type    = string
  default = null
}

variable "open_port" {
  type    = number
  default = null
}

variable "vpc_id" {
  description = "defaults to regional"
  type        = string
  default     = null
}

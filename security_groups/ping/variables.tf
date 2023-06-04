variable "name" {
  type    = string
  default = null
}

variable "open_port" {
  type    = number
  default = 8
}

variable "vpc_id" {
  description = "defaults to regional"
  type        = string
  default     = null
}

variable "ingress_cidr_list" {
    type = list(string)
}

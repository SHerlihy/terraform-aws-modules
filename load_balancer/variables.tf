variable "open_port" {
type = number
default = null
}

variable "vpc_id" {
    type = string
}

variable "subnet_ids" {
  description = "The subnet IDs to deploy to"
  type        = list(string)
}

variable "security_group_ids" {
    type = list(string)
    default = null
}
